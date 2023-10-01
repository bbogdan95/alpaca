package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bbogdan95/alpaca/pkg/engine"
)

func main() {
	engine.InitAll()

	board := &engine.Board{}
	engine.InitHashTable(board, 4)

	s := &engine.SearchInfo{}

	fmt.Println(`
                       ∩~~∩ 
                      ξ ･×･ξ 
                      ξ  ~ ξ 
                      ξ    ξ 
                      ξ   “~～~～〇
                      ξ           ξ	
       d8888 888      ξ ξ ξ~～~ξ  ξ                                
      d88888 888      ξ_ξξ_ξ　ξ_ξξ_ξ                               
     d88P888 888                                     
    d88P 888 888 88888b.   8888b.   .d8888b  8888b.  
   d88P  888 888 888 "88b     "88b d88P"        "88b 
  d88P   888 888 888  888 .d888888 888      .d888888 
 d8888888888 888 888 d88P 888  888 Y88b.    888  888 
d88P     888 888 88888P"  "Y888888  "Y8888P "Y888888 
                 888                                 
                 888                                 
                 888     `)
	fmt.Printf("\nType 'console' for console mode...\n\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		command := strings.TrimSpace(line)

		if line == "\n" {
			continue
		}

		if command == "uci" {
			err := engine.UCILoop(board, s)
			if err != nil {
				panic(err)
			}
			if s.Quit == 1 {
				break
			}

			continue
		}

		if command == "xboard" {
			err := engine.XBoardLoop(board, s)
			if err != nil {
				panic(err)
			}

			if s.Quit == 1 {
				break
			}

			continue
		}

		if command == "console" {
			err := engine.ConsoleLoop(board, s)
			if err != nil {
				panic(err)
			}

			if s.Quit == 1 {
				break
			}

			continue
		}
	}

}
