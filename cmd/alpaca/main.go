package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bbogdan95/alpaca/pkg/engine"
)

var WAC1 = "r1b1k2r/ppppnppp/2n2q2/2b5/3NP3/2P1B3/PP3PPP/RN1QKB1R w KQkq - 0 1"

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
			s.Depth = 6
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
