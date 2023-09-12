package main

import (
	"github.com/bbogdan95/alpaca/pkg/engine"
	"github.com/bbogdan95/alpaca/pkg/perft"
)

func main() {
	engine.InitAll()

	board := &engine.Board{}
	board.ParseFen(engine.PERFTFEN)
	board.CheckBoard()

	perft.PerftTest(3, board, true)
}
