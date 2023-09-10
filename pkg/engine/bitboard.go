package engine

import (
	"fmt"
	"io"
)

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

func PopBit(bb *uint64) int {
	b := *bb ^ (*bb - 1)
	fold := (b & 1) ^ (b >> 32)
	*bb &= (*bb - 1)

	return BitTable[(fold*0x783a9b23)>>26]
}

func CountBits(b uint64) int {
	count := 0
	for b > 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func InitBitMasks() {
	for index := 0; index < 64; index++ {
		SetMask[index] = 1 << index
		ClearMask[index] = ^SetMask[index]
	}
}

func ClearBit(bb *uint64, sq int) {
	*bb &= ClearMask[sq]
}

func SetBit(bb *uint64, sq int) {
	*bb |= SetMask[sq]
}
