package main

import (
	"alpaca-chess/bitboard"
	"fmt"
	"os"
)

/*

r n b q k b n r
p p p p p p p p



P P P P P P P P
R N B Q K B N R

 take this representation anbd transform it into 12 bitboards
 have a custom type for [12]Bitboard
 this type would have a String method that would combine all bitboards and
 print the board in this human readable format

 arrayToBitboards
*/

func main() {
	input := [8][8]string{
		{"r", "n", "b", "q", "k", "b", "n", "r"},
		{"p", "p", "p", "p", "p", "p", "p", "p"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"P", "P", "P", "P", "P", "P", "P", "P"},
		{"R", "N", "B", "Q", "K", "B", "N", "R"},
	}

	bb := bitboard.NewChessBoard(input)

	possibleWhiteMoves("a", bb)

}

const (
	FILE_A          = 72340172838076673
	FILE_H          = 9259542123273814144
	FILE_AB         = 217020518514230019
	FILE_GH         = 13889313184910721216
	RANK_1          = 18374686479671623680
	RANK_4          = 1095216660480
	RANK_5          = 4278190080
	RANK_8          = 255
	CENTER          = 103481868288
	EXTENDED_CENTER = 66229406269440
	KING_SIDE       = 1085102592571150096
	QUEEN_SIDE      = 1085102592571150095
	KING_B7         = 460039
	KNIGHT_C6       = 43234889994
)

var (
	NOT_WHITE_PIECES = bitboard.Bitboard(0)
	BLACK_PIECES     = bitboard.Bitboard(0)
	EMPTY            = bitboard.Bitboard(0)
)

func possibleWhiteMoves(history string, cb *bitboard.ChessBoard) {
	NOT_WHITE_PIECES = cb.WP | cb.WN | cb.WB | cb.WR | cb.WQ | cb.WK | cb.BK
	BLACK_PIECES = cb.BP | cb.BN | cb.BB | cb.BR | cb.BQ
	EMPTY = ^(cb.WP | cb.WN | cb.WN | cb.WR | cb.WQ | cb.WB | cb.WK | cb.BP | cb.BN | cb.BN | cb.BR | cb.BQ | cb.BB | cb.BK)

	fmt.Println("not white pieces")
	NOT_WHITE_PIECES.String(os.Stdout)

	fmt.Println("black pieces")
	BLACK_PIECES.String(os.Stdout)

	fmt.Println("empty")
	EMPTY.String(os.Stdout)

	fmt.Println("white pawns")
	cb.WP.String(os.Stdout)

	WPAttackLeft := (cb.WP >> 7)
	fmt.Println("white pawns attack left")
	WPAttackLeft.String(os.Stdout)

	WPAttackRight := ((cb.WP >> 7) | cb.WP) & ^cb.WP
	fmt.Println("white pawns attack right")
	WPAttackRight.String(os.Stdout)

}
