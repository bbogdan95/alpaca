package main

import (
	"fmt"
	"os"
)

var leafNodes uint64

func Perft(depth int, b *Board) {
	b.CheckBoard()

	if depth == 0 {
		leafNodes++
		return
	}

	ml := &MoveList{}
	GenerateAllMoves(b, ml)

	for i := 0; i < ml.Count; i++ {
		res, err := b.MakeMove(ml.Moves[i].Move)
		if err != nil {
			panic(err)
		}

		if res == FALSE {
			continue
		}

		Perft(depth-1, b)
		b.TakeMove()
	}
}

func PerftTest(depth int, b *Board) {
	b.CheckBoard()

	b.PrintBoard(os.Stdout)
	fmt.Fprintf(os.Stdout, "\nStarting Test To Depth:%d\n", depth)

	leafNodes = 0
	ml := &MoveList{}
	GenerateAllMoves(b, ml)

	for i := 0; i < ml.Count; i++ {
		res, err := b.MakeMove(ml.Moves[i].Move)
		if err != nil {
			panic(err)
		}

		if res == FALSE {
			continue
		}

		cumnodes := leafNodes
		Perft(depth-1, b)
		b.TakeMove()
		oldnodes := leafNodes - cumnodes

		fmt.Fprintf(os.Stdout, "move %d : %s : %1d\n", i+1, PrintMove(ml.Moves[i].Move), oldnodes)
	}

	fmt.Fprintf(os.Stdout, "\nTest Complete : %d nodes visited\n", leafNodes)
}
