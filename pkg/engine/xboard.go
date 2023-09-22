package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func MainLoop() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Alpaca v0.1\n")

	board := &Board{}
	s := &SearchInfo{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if line == "mirrortest\r\n" {
			MirrorEvalTest(board)
		}

		if len(line) >= 8 && line[:8] == "position" {
			board.ParsePosition(line)
		}

		if line == "uci\r\n" {
			err := UCILoop(board, s)
			if err != nil {
				return err
			}
		}

		if line == "quit\r\n" {
			break
		}

	}

	return nil

}

func MirrorEvalTest(b *Board) {
	file, err := os.Open("./mirror.epd")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(-1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	positions := 0

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		b.ParseFen(parts[0])
		positions++

		ev1 := EvalPosition(b)
		b.MirrorBoard()
		ev2 := EvalPosition(b)

		if ev1 != ev2 {
			fmt.Printf("\n\n\n")
			b.ParseFen(parts[0])
			b.PrintBoard(os.Stdout)
			b.MirrorBoard()
			b.PrintBoard(os.Stdout)
			fmt.Printf("\n\nMirror Fail:\n%s\n", parts[0])

			return
		}

		if positions%1000 == 0 {
			fmt.Printf("position %d\n", positions)
		}
	}
}
