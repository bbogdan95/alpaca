package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bbogdan95/alpaca/pkg/engine"
)

var WAC1 = "2rr3k/pp3pp1/1nnqbN1p/3pN3/2pP4/2P3Q1/PPB4P/R4RK1 w - -"

func main() {
	engine.InitAll()

	board := &engine.Board{}
	board.ParseFen(WAC1)
	board.CheckBoard()

	s := &engine.SearchInfo{}

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
		} else if input[0] == 's' {
			s.Depth = 4
			engine.SearchPosition(board, s)

		} else if input[0] == 't' {
			board.TakeMove()
			continue
		} else {
			move, err := engine.ParseMove(input, board)
			if err != nil {
				panic(err)
			}
			if move != engine.NOMOVE {
				board.PvTable.StorePvMove(board, move)
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
