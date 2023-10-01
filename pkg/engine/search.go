package engine

import (
	"fmt"
	"math"
	"os"
	"time"
)

type SearchInfo struct {
	StartTime time.Time
	StopTime  time.Time
	Depth     int
	Depthset  int
	Timeset   int
	Movestogo int
	Infinite  int

	Nodes   uint64
	Quit    int
	Stopped int

	FailHigh      float64
	FailHighFirst float64

	GameMode     int
	PostThinking int
}

// SearchPosition initiates the chess engine's search from the current position on the given board.
// It uses the specified SearchInfo struct to guide the search parameters and store search-related information.
// The function performs an iterative deepening search, updating the principal variation and best move found
// during the search process.
func SearchPosition(b *Board, s *SearchInfo) {
	bestMove := NOMOVE
	ClearForSearch(b, s)

	for currentDepth := 1; currentDepth <= s.Depth; currentDepth++ {
		bestScore := AlphaBeta(-INFINITE, INFINITE, currentDepth, TRUE, b, s)

		if s.Stopped == 1 {
			break
		}

		pvMoves := GetPvLine(currentDepth, b)
		bestMove = b.PvArray[0]

		elapsed := time.Since(s.StartTime).Milliseconds()

		if s.GameMode == UCIMODE {
			fmt.Printf("info score cp %d depth %d nodes %1d time %d ", bestScore, currentDepth, s.Nodes, elapsed)
		} else if s.GameMode == XBOARDMODE && s.PostThinking == TRUE {
			fmt.Printf("%d %d %d %1d ", currentDepth, bestScore, elapsed, s.Nodes)
		} else if s.PostThinking == TRUE {
			fmt.Printf("score:%d depth:%d nodes:%1d time:%d(ms) ", bestScore, currentDepth, s.Nodes, elapsed)
		}

		if s.GameMode == UCIMODE || s.PostThinking == TRUE {
			fmt.Printf("pv")
			for i := 0; i < pvMoves; i++ {
				fmt.Printf(" %s", PrintMove(b.PvArray[i]))
			}
			fmt.Printf("\n")
		}
	}

	if s.GameMode == UCIMODE {
		fmt.Printf("bestmove %s\n", PrintMove(bestMove))
	} else if s.GameMode == XBOARDMODE {
		fmt.Printf("move %s\n", PrintMove(bestMove))
		b.MakeMove(bestMove)
	} else {
		fmt.Printf("\n\n***Alpaca makes move %s***\n\n", PrintMove(bestMove))
		b.MakeMove(bestMove)
		b.PrintBoard(os.Stdout)
	}
}

func ClearForSearch(b *Board, s *SearchInfo) {
	for i := 0; i < 13; i++ {
		for j := 0; j < BRD_SQ_NUM; j++ {
			b.SearchHistory[i][j] = 0
		}
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < MAXDEPTH; j++ {
			b.SearchKillers[i][j] = 0
		}
	}

	b.HashTable.OverWrite = 0
	b.HashTable.Hit = 0
	b.HashTable.Cut = 0

	b.Ply = 0

	s.Stopped = 0
	s.Nodes = 0
	s.FailHighFirst = 0
	s.FailHigh = 0
}

// Performs a quiescence search on the given chess position, evaluating
// capturing moves and other forcing moves to handle tactical threats. It recursively explores
// these moves to a limited depth to accurately evaluate positions where the material balance
// is unstable due to tactical complications.
//
// This function helps avoid the horizon effect, where the engine might miss tactical threats
// just beyond the search horizon, by focusing on capturing and forcing moves that can significantly
// impact the position's evaluation. It uses a limited-depth recursive approach to assess the
// stability of the position and the consequences of potential captures and checks.
func Quiescence(alpha, beta int, b *Board, s *SearchInfo) int {
	b.CheckBoard()
	if beta <= alpha {
		panic("something went wrong")
	}

	if s.Nodes&2047 == 0 {
		CheckUp(s)
	}

	s.Nodes++

	if b.IsRepetition() || b.FiftyMove >= 100 {
		return 0
	}

	if b.Ply > MAXDEPTH-1 {
		return EvalPosition(b)
	}

	score := EvalPosition(b)

	if score >= beta {
		return beta
	}

	if score > alpha {
		alpha = score
	}

	var ml MoveList
	GenerateAllCaptures(b, &ml)

	legal := 0
	score = -INFINITE

	for i := 0; i < ml.Count; i++ {

		PickNextMove(i, &ml)

		res, err := b.MakeMove(ml.Moves[i].Move)
		if err != nil {
			panic(err)
		}
		if res == FALSE {
			continue
		}

		legal++

		score = -Quiescence(-beta, -alpha, b, s)
		b.TakeMove()

		if s.Stopped == TRUE {
			return 0
		}

		if score > alpha {
			if score >= beta {
				if legal == 1 {
					s.FailHighFirst++
				}
				s.FailHigh++
				return beta
			}

			alpha = score
		}
	}

	return alpha
}

func AlphaBeta(alpha, beta, depth, doNull int, b *Board, s *SearchInfo) int {
	b.CheckBoard()

	if depth <= 0 {
		return Quiescence(alpha, beta, b, s)
	}

	if s.Nodes&2047 == 0 {
		CheckUp(s)
	}

	s.Nodes++

	if (b.IsRepetition() || b.FiftyMove >= 100) && b.Ply != 0 {
		return 0
	}

	if b.Ply > MAXDEPTH-1 {
		return EvalPosition(b)
	}

	inCheck := SqAttacked(b.KingSq[b.Side], b.Side^1, b)
	if inCheck == 1 {
		depth++
	}

	score := -INFINITE
	pvMove := NOMOVE

	if ProbeHashEntry(b, &pvMove, &score, alpha, beta, depth) == TRUE {
		b.HashTable.Cut++
		return score
	}

	if doNull == 1 && inCheck == 0 && b.Ply > 0 && b.BigPCE[b.Side] > 0 && depth >= 4 {
		b.MakeNullMove()
		score = -AlphaBeta(-beta, -beta+1, depth-4, FALSE, b, s)
		b.TakeNullMove()
		if s.Stopped == TRUE {
			return 0
		}

		if score >= beta && math.Abs(float64(score)) < ISMATE {
			return beta
		}
	}

	var ml MoveList
	GenerateAllMoves(b, &ml)

	legal := 0
	oldAlpha := alpha
	bestMove := NOMOVE
	bestScore := -INFINITE
	score = -INFINITE

	if pvMove != NOMOVE {
		for i := 0; i < ml.Count; i++ {
			if ml.Moves[i].Move == pvMove {
				ml.Moves[i].Score = 2000000
				break
			}
		}
	}

	for i := 0; i < ml.Count; i++ {

		PickNextMove(i, &ml)

		res, err := b.MakeMove(ml.Moves[i].Move)
		if err != nil {
			panic(err)
		}
		if res == FALSE {
			continue
		}

		legal++

		score = -AlphaBeta(-beta, -alpha, depth-1, TRUE, b, s)
		b.TakeMove()

		if s.Stopped == TRUE {
			return 0
		}

		if score > bestScore {
			bestScore = score
			bestMove = ml.Moves[i].Move

			if score > alpha {
				if score >= beta {
					if legal == 1 {
						s.FailHighFirst++
					}
					s.FailHigh++

					if ml.Moves[i].Move&MoveFlagCapture == 0 {
						b.SearchKillers[1][b.Ply] = b.SearchKillers[0][b.Ply]
						b.SearchKillers[0][b.Ply] = ml.Moves[i].Move
					}

					StoreHashEntry(b, bestMove, beta, HFBETA, depth)

					return beta
				}

				alpha = score

				if ml.Moves[i].Move&MoveFlagCapture == 0 {
					b.SearchHistory[b.Pieces[GetFrom(bestMove)]][GetToSq(bestMove)] += depth
				}
			}
		}

	}

	if legal == 0 {
		if inCheck == 1 {
			return -INFINITE + b.Ply
		} else {
			return 0
		}
	}

	if alpha != oldAlpha {
		StoreHashEntry(b, bestMove, bestScore, HFEXACT, depth)
	} else {
		StoreHashEntry(b, bestMove, alpha, HFALPHA, depth)
	}

	return alpha
}

func PickNextMove(moveNum int, ml *MoveList) {
	var temp Move
	bestScore := 0
	bestNum := moveNum

	for i := moveNum; i < ml.Count; i++ {
		if ml.Moves[i].Score > bestScore {
			bestScore = ml.Moves[i].Score
			bestNum = i
		}
	}

	temp = ml.Moves[moveNum]
	ml.Moves[moveNum] = ml.Moves[bestNum]
	ml.Moves[bestNum] = temp
}

// Check if time is up, or interrupted from GUI
func CheckUp(s *SearchInfo) {
	now := time.Now()
	if s.Timeset == 1 && now.After(s.StopTime) {
		s.Stopped = TRUE
	}
	ReadInput(s)
}
