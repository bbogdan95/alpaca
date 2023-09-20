package engine

import (
	"fmt"
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
}

// Iterative deepening
// for depth = 1 to maxdepth, search with aplha-beta
// if we have enough time, search for depth = 2, and so on
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
		fmt.Printf("info score cp %d depth %d nodes %1d time %d ", bestScore, currentDepth, s.Nodes, elapsed)

		//pvMoves := GetPvLine(currentDepth, b)
		fmt.Printf("pv")
		for i := 0; i < pvMoves; i++ {
			fmt.Printf(" %s", PrintMove(b.PvArray[i]))
		}
		fmt.Printf("\n")
		//fmt.Printf("Ordering: %.2f\n", (s.FailHighFirst / s.FailHigh))
	}

	fmt.Printf("bestmove %s\n", PrintMove(bestMove))
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

	b.PvTable.ClearPvTable()
	b.Ply = 0

	s.Stopped = 0
	s.Nodes = 0
	s.FailHighFirst = 0
	s.FailHigh = 0
}

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
	oldAlpha := alpha
	bestMove := NOMOVE
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
			bestMove = ml.Moves[i].Move
		}
	}

	if alpha != oldAlpha {
		b.PvTable.StorePvMove(b, bestMove)
	}

	return alpha
}

func AlphaBeta(alpha, beta, depth, doNull int, b *Board, s *SearchInfo) int {
	b.CheckBoard()

	if depth == 0 {
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

	var ml MoveList
	GenerateAllMoves(b, &ml)

	legal := 0
	oldAlpha := alpha
	bestMove := NOMOVE
	score := -INFINITE
	PvMove := b.PvTable.ProbePvTable(b)

	if PvMove != NOMOVE {
		for i := 0; i < ml.Count; i++ {
			if ml.Moves[i].Move == PvMove {
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

				return beta
			}

			alpha = score
			bestMove = ml.Moves[i].Move

			if ml.Moves[i].Move&MoveFlagCapture == 0 {
				b.SearchHistory[b.Pieces[GetFrom(bestMove)]][GetToSq(bestMove)] += depth
			}
		}
	}

	if legal == 0 {
		if inCheck == 1 {
			return -MATE + b.Ply
		} else {
			return 0
		}
	}

	if alpha != oldAlpha {
		b.PvTable.StorePvMove(b, bestMove)
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
	//ReadInput(s)
}
