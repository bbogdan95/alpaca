package bitboard

import (
	"fmt"
	"io"
)

type Bitboard uint64
type ChessBoard struct {
	WP Bitboard
	WR Bitboard
	WN Bitboard
	WB Bitboard
	WQ Bitboard
	WK Bitboard
	BP Bitboard
	BR Bitboard
	BN Bitboard
	BB Bitboard
	BQ Bitboard
	BK Bitboard
}

const (
	BoardSize = 8
)

func NewChessBoard(board [8][8]string) *ChessBoard {
	var b = ChessBoard{}

	for file := 0; file < BoardSize; file++ {
		for rank := 0; rank < BoardSize; rank++ {
			square := rank*BoardSize + file

			switch board[rank][file] {
			case "r":
				b.BR |= (1 << uint(square))
			case "n":
				b.BN |= (1 << uint(square))
			case "b":
				b.BB |= (1 << uint(square))
			case "q":
				b.BQ |= (1 << uint(square))
			case "k":
				b.BK |= (1 << uint(square))
			case "p":
				b.BP |= (1 << uint(square))
			case "R":
				b.WR |= (1 << uint(square))
			case "N":
				b.WN |= (1 << uint(square))
			case "B":
				b.WB |= (1 << uint(square))
			case "Q":
				b.WQ |= (1 << uint(square))
			case "K":
				b.WK |= (1 << uint(square))
			case "P":
				b.WP |= (1 << uint(square))
			}
		}
	}

	return &b
}

func (cb *ChessBoard) String(output io.Writer) {
	out := [8][8]string{}

	for rank := 0; rank < BoardSize; rank++ {
		for file := 0; file < BoardSize; file++ {
			square := rank*BoardSize + file
			if cb.BB&(1<<uint(square)) != 0 {
				out[rank][file] = "b"
			}
			if cb.BK&(1<<uint(square)) != 0 {
				out[rank][file] = "k"
			}
			if cb.BQ&(1<<uint(square)) != 0 {
				out[rank][file] = "q"
			}
			if cb.BR&(1<<uint(square)) != 0 {
				out[rank][file] = "r"
			}
			if cb.BN&(1<<uint(square)) != 0 {
				out[rank][file] = "n"
			}
			if cb.BP&(1<<uint(square)) != 0 {
				out[rank][file] = "p"
			}

			if cb.WB&(1<<uint(square)) != 0 {
				out[rank][file] = "B"
			}
			if cb.WK&(1<<uint(square)) != 0 {
				out[rank][file] = "K"
			}
			if cb.WQ&(1<<uint(square)) != 0 {
				out[rank][file] = "Q"
			}
			if cb.WR&(1<<uint(square)) != 0 {
				out[rank][file] = "R"
			}
			if cb.WN&(1<<uint(square)) != 0 {
				out[rank][file] = "N"
			}
			if cb.WP&(1<<uint(square)) != 0 {
				out[rank][file] = "P"
			}
		}
	}

	for rank := 0; rank < BoardSize; rank++ {
		for file := 0; file < BoardSize; file++ {
			fmt.Fprintf(output, "%s ", out[rank][file])
		}
		fmt.Fprintln(output, "")
	}
}

func (bb *Bitboard) String(output io.Writer) {
	for rank := 0; rank < BoardSize; rank++ {
		for file := 0; file < BoardSize; file++ {
			square := rank*BoardSize + file
			if *bb&(1<<uint(square)) != 0 {
				fmt.Fprint(output, "1 ")
			} else {
				fmt.Fprint(output, "0 ")
			}
		}
		fmt.Fprintln(output)
	}
}

func (bb *Bitboard) SetBit(square int) *Bitboard {
	*bb |= (1 << uint(square))
	return bb
}

func SquareToIndex(square string) (int, error) {
	if len(square) != 2 {
		return 0, fmt.Errorf("invalid square notation: %s", square)
	}

	file := int(square[0] - 'A')
	rank := 7 - int(square[1]-'1')

	if file < 0 || file >= BoardSize || rank < 0 || rank >= BoardSize {
		return 0, fmt.Errorf("invalid square notation: %s", square)
	}

	return rank*BoardSize + file, nil
}

func (bb *Bitboard) SetBitBySquare(square string) error {
	index, err := SquareToIndex(square)
	if err != nil {
		return err
	}

	bb.SetBit(index)
	return nil
}
