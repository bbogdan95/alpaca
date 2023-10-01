package engine

import (
	"fmt"
)

// KnightDir represents the possible move directions for a knight on the chessboard.
// The knight can move in an L-shape: two squares in one direction and one square perpendicular to that direction.
// The array contains relative square indices representing these knight move directions.
var KnightDir = [8]int{-8, -19, -21, -12, 8, 19, 21, 12}

// RookDir represents the possible move directions for a rook on the chessboard.
// Rooks can move vertically (along files) or horizontally (along ranks).
// The array contains relative square indices representing these rook move directions.
var RookDir = [4]int{-1, -10, 1, 10}

// BishopDir represents the possible move directions for a bishop on the chessboard.
// Bishops can move diagonally.
// The array contains relative square indices representing these bishop move directions.
var BishopDir = [4]int{-9, -11, 11, 9}

// KingDir represents the possible move directions for a king on the chessboard.
// Kings can move one square in any direction: horizontally, vertically, or diagonally.
// The array contains relative square indices representing these king move directions.
var KingDir = [8]int{-1, -10, 1, 10, -9, -11, 11, 9}

// SqAttacked determines if a specific square on the chessboard is attacked by a given side.
//
// Parameters:
//
//	sq: The square index to check for an attack.
//	side: The side for which the attack is being checked. Use constants WHITE or BLACK.
//	b: A pointer to the chessboard struct representing the current board state.
//
// Returns:
//
//	TRUE if the square is attacked by the specified side, FALSE otherwise.
//
// Note:
//
//	The function handles pawns, knights, rooks, bishops, queens, and kings to determine attacks.
//	It performs boundary checks and ensures the square and side are valid before processing.
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

// ShowSqAttackedBySide prints the squares attacked by the specified side on the given chessboard.
//
// Parameters:
//
//	side: The side for which the attacked squares are to be displayed. Use constants WHITE or BLACK.
//	b: A pointer to the chessboard for which the attacked squares are to be calculated.
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
