package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func XBoardLoop(b *Board, s *SearchInfo) error {
	reader := bufio.NewReader(os.Stdin)

	s.GameMode = XBOARDMODE
	s.PostThinking = TRUE
	depth := -1
	movestogo := [2]int{30, 30}
	movetime := -1
	inc := 0
	engineSide := BOTH
	t := 0
	mps := 0
	inBuf := ""
	command := ""
	MB := 64
	timeLeft := 0

	for {
		if b.Side == engineSide && CheckResult(b) == FALSE {
			s.StartTime = time.Now()
			s.Depth = depth

			if t != -1 {
				s.Timeset = TRUE
				t /= movestogo[b.Side]
				t -= 50
				durationOffset := time.Duration(t+inc) * time.Millisecond
				s.StopTime = s.StartTime.Add(durationOffset)
			}

			if depth == -1 || depth > MAXDEPTH {
				s.Depth = MAXDEPTH
			}

			fmt.Printf("time:%d start:%s stop:%s depth:%d timeset:%d movestgoto:%d mps:%d\n", t, s.StartTime, s.StopTime, s.Depth, s.Timeset, movestogo[b.Side], mps)
			SearchPosition(b, s)

			if mps != 0 {
				movestogo[b.Side^1]--
				if movestogo[b.Side^1] < 1 {
					movestogo[b.Side^1] = mps
				}
			}
		}

		inBuf, _ = reader.ReadString('\n')
		inBuf = strings.TrimSpace(inBuf)

		if len(inBuf) == 0 {
			continue
		}

		command = strings.Fields(inBuf)[0]

		fmt.Printf("command seen:%s\n", inBuf)

		switch command {
		case "quit":
			s.Quit = TRUE
			return nil
		case "force":
			engineSide = BOTH
		case "protover":
			PrintOptions()
		case "sd":
			if n, err := fmt.Sscanf(inBuf, "sd %d", &depth); err == nil && n == 1 {
				fmt.Printf("DEBUG depth:%d\n", depth)
			}
		case "st":
			if n, err := fmt.Sscanf(inBuf, "st %d", &movetime); err == nil && n == 1 {
				fmt.Printf("DEBUG movetime:%d\n", movetime)
			}
		case "time":
			if n, err := fmt.Sscanf(inBuf, "time %d", &t); err == nil && n == 1 {
				t *= 10
				fmt.Printf("DEBUG time:%d\n", t)
			}
		case "memory":
			n, err := fmt.Sscanf(inBuf, "memory %d", &MB)
			if err != nil {
				panic(err)
			}
			if n != 1 {
				panic("error reading Hash value")
			}
			if MB < 4 {
				MB = 4
			}
			if MB > 2048 {
				MB = 2048
			}

			fmt.Printf("Set Hash to %d MB\n", MB)
			InitHashTable(b, MB)
		case "level":
			movetime = -1
			sec := 0
			if n, _ := fmt.Sscanf(inBuf, "level %d %d %d", &mps, &timeLeft, &inc); n != 3 {
				fmt.Sscanf(inBuf, "level %d %d:%d %d", &mps, &timeLeft, &sec, &inc)
				fmt.Printf("DEBUG level with :\n")
			} else {
				fmt.Printf("DEBUG level without :\n")
			}

			timeLeft *= 60000
			timeLeft += sec * 1000
			movestogo[0] = 30
			movestogo[1] = 30

			if mps != 0 {
				movestogo[0] = mps
				movestogo[1] = mps
			}

			t = -1
			fmt.Printf("DEBUG level timeLeft:%d movesToGo:%d inc:%d mps:%d\n", timeLeft, movestogo[0], inc, mps)
		case "ping":
			fmt.Printf("pong%s\n", inBuf[4:])
		case "new":
			ClearHashTable(b)
			engineSide = BLACK
			b.ParseFen(START_FEN)
			depth = -1
			t = -1
		case "setboard":
			engineSide = BOTH
			b.ParseFen(strings.TrimSpace(strings.TrimPrefix(inBuf, "setboard")))
		case "go":
			engineSide = b.Side
		case "usermove":
			movestogo[b.Side]--
			move, _ := ParseMove(strings.TrimSpace(strings.TrimPrefix(inBuf, "usermove")), b)
			if move != NOMOVE {
				b.MakeMove(move)
				b.Ply = 0
			}
		}
	}
}

func PrintOptions() {
	fmt.Println("feature ping=1 setboard=1 colors=0 usermove=1 memory=1")
	fmt.Println("feature done=1")
}

func ConsoleLoop(b *Board, s *SearchInfo) error {
	fmt.Println("Alpaca - Console Mode")
	fmt.Println("Type help for commands")

	s.GameMode = CONSOLEMODE
	s.PostThinking = TRUE
	reader := bufio.NewReader(os.Stdin)

	depth := MAXDEPTH
	movetime := 3000
	engineSide := BOTH
	inBuf := ""
	command := ""

	for {
		if b.Side == engineSide && CheckResult(b) == FALSE {
			s.StartTime = time.Now()
			s.Depth = depth

			if movetime != 0 {
				s.Timeset = TRUE
				s.StopTime = s.StartTime.Add(time.Duration(movetime) * time.Millisecond)
			}

			SearchPosition(b, s)
		}

		fmt.Print("\nAlpaca > ")

		inBuf, _ = reader.ReadString('\n')
		inBuf = strings.TrimSpace(inBuf)

		if len(inBuf) == 0 {
			continue
		}

		command = strings.Fields(inBuf)[0]

		switch command {
		case "help":
			fmt.Println("Commands:")
			fmt.Println("quit - quit game")
			fmt.Println("force - computer will not think")
			fmt.Println("print - show board")
			fmt.Println("post - show thinking")
			fmt.Println("nopost - do not show thinking")
			fmt.Println("new - start a new game")
			fmt.Println("go - set computer thinking")
			fmt.Println("depth x - set depth to x")
			fmt.Println("time x - set thinking time to x seconds (depth still applies if set)")
			fmt.Println("view - show current depth and movetime settings")
			fmt.Println("setboard x - set position to fen x")
			fmt.Println("** note ** - to reset time and depth, set to 0")
			fmt.Println("enter moves using b7b8q notation")
		case "mirror":
			engineSide = BOTH
			MirrorEvalTest(b)
		case "eval":
			b.PrintBoard(os.Stdout)
			fmt.Printf("Eval:%d\n", EvalPosition(b))
			b.MirrorBoard()
			b.PrintBoard(os.Stdout)
			fmt.Printf("Eval:%d\n", EvalPosition(b))
		case "setboard":
			engineSide = BOTH
			fen := strings.TrimSpace(strings.TrimPrefix(inBuf, "setboard"))
			b.ParseFen(fen)
		case "quit":
			s.Quit = TRUE
			return nil
		case "post":
			s.PostThinking = TRUE
		case "print":
			b.PrintBoard(os.Stdout)
			continue
		case "nopost":
			s.PostThinking = FALSE
		case "force":
			engineSide = BOTH
		case "view":
			if depth == MAXDEPTH {
				fmt.Printf("depth not set ")
			} else {
				fmt.Printf("depth %d", depth)
			}

			if movetime != 0 {
				fmt.Printf(" movetime %ds\n", movetime/1000)
			} else {
				fmt.Println(" movetime not set")
			}
		case "depth":
			var d int
			n, err := fmt.Sscanf(inBuf, "depth %d", &d)
			if err == nil && n == 1 {
				depth = d
				if depth == 0 {
					depth = MAXDEPTH
				}
			}
		case "time":
			var t int
			n, err := fmt.Sscanf(inBuf, "time %d", &t)
			if err == nil && n == 1 {
				movetime = t * 1000
			}
		case "new":
			ClearHashTable(b)
			engineSide = BLACK
			b.ParseFen(START_FEN)
		case "go":
			engineSide = b.Side
		default:
			move, err := ParseMove(inBuf, b)
			if err != nil || move == NOMOVE {
				fmt.Printf("Command unknown: %s\n", inBuf)
			} else {
				b.MakeMove(move)
				b.Ply = 0
			}
		}
	}
}

func CheckResult(b *Board) int {
	if b.FiftyMove > 100 {
		fmt.Printf("1/2-1/2 {fifty move rule (claimed by Alpaca)}\n")
		return 1
	}

	if ThreeFoldRepetition(b) >= 2 {
		fmt.Printf("1/2-1/2 {3-fold repetition (claimed by Alpaca)}\n")
		return 1
	}

	if DrawMaterial(b) == 1 {
		fmt.Printf("1/2-1/2 {insufficient material (claimed by Alpaca)}\n")
		return 1
	}

	var ml MoveList
	GenerateAllMoves(b, &ml)
	found := 0

	for i := 0; i < ml.Count; i++ {
		res, _ := b.MakeMove(ml.Moves[i].Move)
		if res == 0 {
			continue
		}
		found++
		b.TakeMove()
		break
	}

	if found != 0 {
		return 0
	}

	inCheck := SqAttacked(b.KingSq[b.Side], b.Side^1, b)
	if inCheck == TRUE {
		if b.Side == WHITE {
			fmt.Printf("0-1 {black mates (claimed by Alpaca)}\n")
			return 1
		} else {
			fmt.Printf("0-1 {white mates (claimed by Alpaca)}\n")
			return 1
		}
	} else {
		fmt.Printf("\n1/2-1/2 {stalemate (claimed by Alpaca)}\n")
		return 1
	}
}

func ThreeFoldRepetition(b *Board) int {
	r := 0

	for i := 0; i < b.HisPly; i++ {
		if b.History[i].PosKey == b.PosKey {
			r++
		}
	}

	return r
}

func DrawMaterial(b *Board) int {
	if b.PCENum[WP] != 0 || b.PCENum[BP] != 0 {
		return 0
	}

	if b.PCENum[WQ] != 0 || b.PCENum[BQ] != 0 || b.PCENum[WR] != 0 || b.PCENum[BR] != 0 {
		return 0
	}

	if b.PCENum[WB] > 1 || b.PCENum[BB] > 1 {
		return 0
	}

	if b.PCENum[WN] > 1 || b.PCENum[BN] > 1 {
		return 0
	}

	if b.PCENum[WN] != 0 && b.PCENum[WB] != 0 {
		return 0
	}

	if b.PCENum[BN] != 0 && b.PCENum[BB] != 0 {
		return 0
	}

	return 1
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
