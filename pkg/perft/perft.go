package perft

import (
	"fmt"
	"os"
	"time"

	"github.com/bbogdan95/alpaca/pkg/engine"
)

var leafNodes uint64

func Perft(depth int, b *engine.Board) error {
	b.CheckBoard()

	if depth == 0 {
		leafNodes++
		return nil
	}

	ml := &engine.MoveList{}
	engine.GenerateAllMoves(b, ml)

	for i := 0; i < ml.Count; i++ {
		res, err := b.MakeMove(ml.Moves[i].Move)
		if err != nil {
			return err
		}

		if res == 0 {
			continue
		}

		err = Perft(depth-1, b)
		if err != nil {
			return err
		}
		b.TakeMove()
	}

	return nil
}

func PerftTest(depth int, b *engine.Board, log bool) (uint64, error) {
	b.CheckBoard()

	if log {
		b.PrintBoard(os.Stdout)
		fmt.Fprintf(os.Stdout, "\nStarting Test To Depth:%d\n", depth)
	}
	start := time.Now()

	leafNodes = 0
	ml := &engine.MoveList{}
	engine.GenerateAllMoves(b, ml)

	for i := 0; i < ml.Count; i++ {
		res, err := b.MakeMove(ml.Moves[i].Move)
		if err != nil {
			return 0, err
		}

		if res == 0 {
			continue
		}

		cumnodes := leafNodes
		err = Perft(depth-1, b)
		if err != nil {
			return 0, err
		}

		b.TakeMove()
		oldnodes := leafNodes - cumnodes

		if log {
			fmt.Fprintf(os.Stdout, "move %d : %s : %1d\n", i+1, engine.PrintMove(ml.Moves[i].Move), oldnodes)
		}
	}

	if log {
		fmt.Fprintf(os.Stdout, "\nTest Complete : %d nodes visited in %s\n", leafNodes, time.Since(start))
	}

	return leafNodes, nil
}
