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

	ClearForSearch(b, s)

	for currentDepth := 1; currentDepth <= s.Depth; currentDepth++ {
		bestScore := AlphaBeta(-INFINITE, INFINITE, currentDepth, TRUE, b, s)
		GetPvLine(currentDepth, b)
		bestMove := b.PvArray[0]

		fmt.Printf("Depth: %d score: %d move: %s nodes: %1d", currentDepth, bestScore, PrintMove(bestMove), s.Nodes)

		pvMoves := GetPvLine(currentDepth, b)
		fmt.Printf("pv")
		for i := 0; i < pvMoves; i++ {
			fmt.Printf(" %s", PrintMove(b.PvArray[i]))
		}
		fmt.Printf("\n")
		fmt.Printf("Ordering: %.2f\n", (s.FailHighFirst / s.FailHigh))
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

	b.PvTable.ClearPvTable()
	b.Ply = 0

	s.StartTime = time.Now()
	s.Stopped = 0
	s.Nodes = 0
	s.FailHighFirst = 0
	s.FailHigh = 0
}

func Quiescence(alpha, beta int, b *Board, s *SearchInfo) int {
	return 0
}

func AlphaBeta(alpha, beta, depth, doNull int, b *Board, s *SearchInfo) int {
	b.CheckBoard()

	if depth == 0 {
		s.Nodes++
		return EvalPosition(b)
	}

	s.Nodes++

	if b.IsRepetition() || b.FiftyMove >= 100 {
		return 0
	}

	if b.Ply > MAXDEPTH-1 {
		return EvalPosition(b)
	}

	var ml MoveList
	GenerateAllMoves(b, &ml)

	legal := 0
	oldAlpha := alpha
	bestMove := NOMOVE
	score := -INFINITE

	for i := 0; i < ml.Count; i++ {
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

	if legal == 0 {
		if SqAttacked(b.KingSq[b.Side], b.Side^1, b) == TRUE {
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
