package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bbogdan95/alpaca/pkg/engine"
	"github.com/bbogdan95/alpaca/pkg/perft"
)

func main() {
	engine.InitAll()

	board := &engine.Board{}
	board.ParseFen(engine.START_FEN)
	board.CheckBoard()

	reader := bufio.NewReader(os.Stdin)

	for {
		board.PrintBoard(os.Stdout)
		fmt.Println("Please enter a move")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if len(input) > 6 {
			panic("invalid move")
		}

		if input[0] == 'q' {
			break
		} else if input[0] == 'p' {
			perft.PerftTest(4, board, true)
		} else if input[0] == 't' {
			board.TakeMove()
			continue
		} else {
			move, err := engine.ParseMove(input, board)
			if err != nil {
				panic(err)
			}
			if move != engine.NOMOVE {
				board.MakeMove(move)
				// if board.IsRepetition() {
				// 	fmt.Printf("REP SEEN\n")
				// }
			} else {
				fmt.Println("Move not parsed: ", input)
			}
		}
	}

}
