package main

import (
	"log"

	"github.com/bbogdan95/alpaca/pkg/engine"
)

var WAC1 = "r1b1k2r/ppppnppp/2n2q2/2b5/3NP3/2P1B3/PP3PPP/RN1QKB1R w KQkq - 0 1"

func main() {
	engine.InitAll()

	err := engine.UCILoop()
	if err != nil {
		log.Println(err)
	}

}
