package main

import (
	engine "alpaca-chess/pkg/engine"
	"alpaca-chess/pkg/perft"
)

func main() {
	engine.InitAll()

	board := &engine.Board{}
	board.ParseFen(engine.PERFTFEN)
	board.CheckBoard()

	perft.PerftTest(4, board)

}
