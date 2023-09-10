package engine

import (
	"fmt"
)

var KnightDir = [8]int{-8, -19, -21, -12, 8, 19, 21, 12}
var RookDir = [4]int{-1, -10, 1, 10}
var BishopDir = [4]int{-9, -11, 11, 9}
var KingDir = [8]int{-1, -10, 1, 10, -9, -11, 11, 9}

func SqAttacked(sq int, side int, b *Board) int {

	if !SqOnBoard(sq) || !SideValid(side) {
		panic("err")
	}

	b.CheckBoard()

	// pawns
	if side == WHITE {
		if b.Pieces[sq-11] == WP || b.Pieces[sq-9] == WP {
			return TRUE
		}
	} else {
		if b.Pieces[sq+11] == BP || b.Pieces[sq+9] == BP {
			return TRUE
		}
	}

	// knights
	for i := 0; i < 8; i++ {
		piece := b.Pieces[sq+KnightDir[i]]
		if piece != OFFBOARD && piece != EMPTY && PieceKnight[piece] == TRUE && PieceCol[piece] == side {
			return TRUE
		}
	}

	// rooks, queens
	for i := 0; i < 4; i++ {
		dir := RookDir[i]
		tSq := sq + dir
		piece := b.Pieces[tSq]

		for piece != OFFBOARD {
			if piece != EMPTY {
				if PieceRookQueen[piece] == TRUE && PieceCol[piece] == side {
					return TRUE
				}
				break
			}
			tSq += dir
			piece = b.Pieces[tSq]
		}
	}

	// bishops, queen
	for i := 0; i < 4; i++ {
		dir := BishopDir[i]
		tSq := sq + dir
		piece := b.Pieces[tSq]

		for piece != OFFBOARD {
			if piece != EMPTY {
				if PieceBishopQueen[piece] == TRUE && PieceCol[piece] == side {
					return TRUE
				}
				break
			}
			tSq += dir
			piece = b.Pieces[tSq]
		}
	}

	// kings
	for i := 0; i < 8; i++ {
		piece := b.Pieces[sq+KingDir[i]]
		if piece != OFFBOARD && piece != EMPTY && PieceKing[piece] == TRUE && PieceCol[piece] == side {
			return TRUE
		}
	}

	return FALSE
}

func ShowSqAttackedBySide(side int, b *Board) {
	fmt.Printf("\n\nSquares attacked by: %c\n", SideChar[side])

	for rank := RANK_8; rank >= RANK_1; rank-- {
		for file := FILE_A; file <= FILE_H; file++ {
			sq := FR2SQ(file, rank)
			if SqAttacked(sq, side, b) == TRUE {
				fmt.Printf(" X ")
			} else {
				fmt.Printf(" - ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n")
}
