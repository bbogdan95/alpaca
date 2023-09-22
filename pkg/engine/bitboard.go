package engine

import (
	"fmt"
	"io"
)

// SetMask is an array of 64 uint64 values, each containing a single bit set at a specific index.
// It is used to set individual bits in a uint64 bitboard.
var SetMask [64]uint64

// ClearMask is an array of 64 uint64 values, each containing all bits set except for a single bit at a specific index.
// It is used to clear individual bits in a uint64 bitboard.
var ClearMask [64]uint64

// BitTable is an array that maps a 6-bit value (0-63) to the index of the least significant set bit (LS1B) in binary representation.
// It is used for efficient bit manipulation operations like finding the LS1B index.
var BitTable = []int{
	63, 30, 3, 32, 25, 41, 22, 33, 15, 50, 42, 13, 11, 53, 19, 34, 61, 29, 2, 51, 21, 43, 45, 10, 18, 47, 1, 54, 9, 57, 0, 35, 62, 31, 40, 4, 49, 5, 52, 26, 60, 6, 23, 44, 46, 27, 56, 16, 7, 39, 48, 24, 59, 14, 12, 55, 38, 28, 58, 20, 37, 17, 36, 8,
}

// PrintBitboard prints a human-readable representation of a bitboard to the specified output writer.
// 'X' is used to represent set bits, and '-' is used to represent cleared bits.
func PrintBitboard(out io.Writer, bitboard uint64) {
	var shiftMe uint64 = 1

	sq := 0
	sq64 := 0

	for rank := RANK_8; rank >= RANK_1; rank-- {
		for file := FILE_A; file <= FILE_H; file++ {
			sq = FR2SQ(file, rank)
			sq64 = SQ64[sq]

			if (shiftMe<<sq64)&bitboard != 0 {
				fmt.Fprintf(out, "%5s", "X")
			} else {
				fmt.Fprintf(out, "%5s", "-")
			}
		}
		fmt.Fprintf(out, "\n")
	}
	fmt.Fprintf(out, "\n\n")
}

func (b *Board) MirrorBoard() {
	tempPieces := [64]int{}
	tempSide := b.Side ^ 1
	tempCastlePerm := 0
	tempEnPassant := NO_SQ

	swapPiece := [13]int{EMPTY, BP, BN, BB, BR, BQ, BK, WP, WN, WB, WR, WQ, WK}

	if b.CastlePerm&WKCA != 0 {
		tempCastlePerm |= BKCA
	}

	if b.CastlePerm&WQCA != 0 {
		tempCastlePerm |= BQCA
	}

	if b.CastlePerm&BKCA != 0 {
		tempCastlePerm |= WKCA
	}

	if b.CastlePerm&BQCA != 0 {
		tempCastlePerm |= WQCA
	}

	if b.EnPassant != NO_SQ {
		tempEnPassant = SQ120[Mirror64[SQ64[b.EnPassant]]]
	}

	for sq := 0; sq < 64; sq++ {
		tempPieces[sq] = b.Pieces[SQ120[Mirror64[sq]]]
	}

	b.ResetBoard()

	for sq := 0; sq < 64; sq++ {
		tp := swapPiece[tempPieces[sq]]
		b.Pieces[SQ120[sq]] = tp
	}

	b.Side = tempSide
	b.CastlePerm = tempCastlePerm
	b.EnPassant = tempEnPassant

	b.PosKey = GeneratePosKey(b)

	UpdateListsMaterial(b)

	b.CheckBoard()
}

// PopBit finds and clears the least significant set bit (LS1B) in a uint64 bitboard.
// It returns the index (0-63) of the LS1B that was cleared.
// This function is used for efficiently finding and removing individual set bits from a bitboard.
// The BitTable array is used to determine the index of the cleared LS1B.
func PopBit(bb *uint64) int {
	b := *bb ^ (*bb - 1)
	fold := (b & 1) ^ (b >> 32)
	*bb &= (*bb - 1)

	return BitTable[(fold*0x783a9b23)>>26]
}

// CountBits counts the number of set (1) bits in a uint64 bitboard.
func CountBits(b uint64) int {
	count := 0
	for b > 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

// InitBitMasks initializes the SetMask and ClearMask arrays to facilitate bit manipulation.
// SetMask contains individual bits set at specific indices, and ClearMask contains bits cleared at specific indices.
func InitBitMasks() {
	for index := 0; index < 64; index++ {
		SetMask[index] = 1 << index
		ClearMask[index] = ^SetMask[index]
	}
}

// ClearBit clears a specific bit in a uint64 bitboard using the ClearMask array.
func ClearBit(bb *uint64, sq int) {
	*bb &= ClearMask[sq]
}

// SetBit sets a specific bit in a uint64 bitboard using the SetMask array.
func SetBit(bb *uint64, sq int) {
	*bb |= SetMask[sq]
}
