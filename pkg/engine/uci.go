package engine

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const INPUTBUFFER = 400 * 6
const NAME = "ALPACAv0.1"

// go depth 6 wtime 180000 btime 100000 binc 1000 winc 1000 movetime 1000 movestogo 40
func (b *Board) ParseGo(line string, s *SearchInfo) error {
	depth := -1
	movestogo := 30
	movetime := -1
	t := -1
	inc := 0
	s.Timeset = FALSE

	parts := strings.Fields(line)
	for i := 1; i < len(parts); i++ {
		switch parts[i] {
		case "depth":
			i++
			depth, _ = strconv.Atoi(parts[i])
		case "movetime":
			i++
			movetime, _ = strconv.Atoi(parts[i])
		case "wtime":
			i++
			t, _ = strconv.Atoi(parts[i])
		case "btime":
			i++
			t, _ = strconv.Atoi(parts[i])
		case "winc":
			i++
			inc, _ = strconv.Atoi(parts[i])
		case "binc":
			i++
			inc, _ = strconv.Atoi(parts[i])
		}
	}

	if movetime != -1 {
		t = movetime
		movestogo = 1
	}

	s.StartTime = time.Now()
	if t != -1 {
		s.Timeset = TRUE
		t /= movestogo
		t -= 50
		to := time.Millisecond * time.Duration(t+inc)
		s.StopTime = s.StartTime.Add(to)
	}

	if depth == -1 {
		s.Depth = MAXDEPTH
	} else {
		s.Depth = depth
	}

	fmt.Printf("time:%d start:%s stop:%s depth:%d timeset:%d\n", t, s.StartTime, s.StopTime, s.Depth, s.Timeset)

	SearchPosition(b, s)

	return nil
}

func (b *Board) ParsePosition(line string) {
	if line[9:17] == "startpos" {
		b.ParseFen(START_FEN)
	} else if line[9:12] == "fen" {
		b.ParseFen(line[13:])
	} else {
		b.ParseFen(START_FEN)
	}

	parts := strings.Split(line, "moves ")
	if len(parts) > 1 && len(parts[1]) > 0 {
		moves := strings.Split(parts[1], " ")

		for _, m := range moves {
			move, _ := ParseMove(m, b)
			if move == NOMOVE {
				break
			}

			b.MakeMove(move)
			b.Ply = 0
		}
	}

	b.PrintBoard(os.Stdout)
}

func UCILoop(board *Board, s *SearchInfo) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("id name %s\n", NAME)
	fmt.Printf("id author Mid\n")
	fmt.Printf("option name Hash type spin default 64 min 4 max 2048\n")
	fmt.Printf("uciok\n")

	MB := 64

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if len(line) >= 7 && line[:7] == "isready" {
			fmt.Printf("readyok\n")
		} else if len(line) >= 8 && line[:8] == "position" {
			board.ParsePosition(line)
		} else if len(line) >= 10 && line[:10] == "ucinewgame" {
			board.ParsePosition("position startpos\n")
		} else if len(line) >= 2 && line[:2] == "go" {
			board.ParseGo(line, s)
		} else if len(line) >= 4 && line[:4] == "quit" {
			s.Quit = TRUE
		} else if len(line) >= 3 && line[:3] == "uci" {
			fmt.Printf("id name %s\n", NAME)
			fmt.Printf("id author Mid\n")
			fmt.Printf("uciok\n")
		} else if len(line) >= 26 && line[:26] == "setoption name Hash value " {
			fmt.Sscanf(line, "%*s %*s %*s %d", &MB)
			if MB < 4 {
				MB = 4
			}
			if MB > 2048 {
				MB = 2048
			}

			fmt.Printf("Set Hash to %d MB\n", MB)
			board.HashTable.ClearHashTable()
		}

		if s.Quit == TRUE {
			break
		}
	}

	return nil
}
