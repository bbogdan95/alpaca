package main

import (
	"alpaca-chess/pkg/engine"
	"alpaca-chess/pkg/perft"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	engine.InitAll()

	PerftTestSuite("./perftsuite.edp")
}

func PerftTestSuite(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(-1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	board := &engine.Board{}

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, ";")

		board.ParseFen(parts[0])
		board.CheckBoard()

		fmt.Printf("%s\n", parts[0])
		for _, test := range parts[1:] {
			// if partsIndex <= 2 {
			testParts := strings.Split(test, " ")

			depthChar := testParts[0][1]
			depthInt := int(depthChar - '0')

			fmt.Printf(" - depth %d - %s - ", depthInt, testParts[1])

			leafNodesCheck, err := strconv.ParseUint(testParts[1], 10, 64)
			if err != nil {
				panic(err)
			}

			leafNodes, err := perft.PerftTest(depthInt, board, false)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%d", leafNodes)

			if leafNodes != leafNodesCheck {
				fmt.Printf("%s\n", "❌")
				os.Exit(-1)
			} else {
				fmt.Printf("%s\n", "✅")
			}
			// }
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
