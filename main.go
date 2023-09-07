package main

import (
	"bufio"
	"fmt"
	"os"
)

var SQ64 [BRD_SQ_NUM]int
var SQ120 [64]int
var SetMask [64]uint64
var ClearMask [64]uint64

var FilesBrd [BRD_SQ_NUM]int
var RanksBrd [BRD_SQ_NUM]int

var PceChar = ".PNBRQKpnbrqk"
var SideChar = "wb-"
var RankChar = "12345678"
var FileChar = "abcdefgh"

var PieceBig = [13]int{FALSE, FALSE, TRUE, TRUE, TRUE, TRUE, TRUE, FALSE, TRUE, TRUE, TRUE, TRUE, TRUE}
var PieceMaj = [13]int{FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, TRUE, TRUE}
var PieceMin = [13]int{FALSE, FALSE, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, FALSE, FALSE, FALSE}
var PieceVal = [13]int{0, 100, 325, 325, 550, 1000, 50000, 100, 325, 325, 550, 1000, 50000}
var PieceCol = [13]int{BOTH, WHITE, WHITE, WHITE, WHITE, WHITE, WHITE, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK}

var PiecePawn = [13]int{FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE}
var PieceKnight = [13]int{FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE}
var PieceKing = [13]int{FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE}
var PieceRookQueen = [13]int{FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, FALSE}
var PieceBishopQueen = [13]int{FALSE, FALSE, FALSE, TRUE, FALSE, TRUE, FALSE, FALSE, FALSE, TRUE, FALSE, TRUE, FALSE}
var PieceSlides = [13]int{FALSE, FALSE, FALSE, TRUE, TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, TRUE, TRUE, FALSE}

// used for movegen
var LoopSlidePieces = [8]int{WB, WR, WQ, 0, BB, BR, BQ, 0}
var LoopSlideIndex = [2]int{0, 4}
var LoopNonSlidePieces = [6]int{WN, WK, 0, BN, BK, 0}
var LoopNonSlideIndex = [2]int{0, 3}
var PieceDir = [13][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
}
var NumDir = [13]int{0, 0, 8, 4, 4, 8, 8, 0, 8, 4, 4, 8, 8}

// Every time we move a piece, we will do castle_permissions &= CastlePerm[from]
// and castle_permissions &= CastlePerm[from]. The result of these operations is 1111 == 15
// except for A1, E1, H1 & A8, E8, H8
// When the rooks of the queen moves on either side, it takes out the castle permissions for that side
// eq. Black queen moves from E8 to E7. castle_permissions &= 3 -> gives 0011 -> which means BLACK side lost castling permissions
// on both queen and king side.
var CastlePerm = [120]int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 7, 15, 15, 15, 3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

var FEN0 = "8/3q4/8/8/4Q3/8/8/8 w - - 0 2"
var FEN1 = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var FEN2 = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
var FEN3 = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
var FEN4 = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

var PAWNMOVES_W = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
var PAWNMOVES_B = "rnbqkbnr/p1p1p3/3p3p/1p1p4/2P1Pp2/8/PP1P1PpP/RNBQKB1R b KQkq e3 0 1"
var KNIGHTSKINGSFEN = "5k2/1n6/4n3/6N1/8/3N4/8/5K2 w - - 0 1"
var ROOKSFEN = "6k1/8/5r2/8/1nR5/5N2/8/6K1 b - - 0 1"
var QUEENSFEN = "6k1/8/4nq2/8/1nQ5/5N2/1N6/6K1 b - - 0 1"
var BISHOPSFEN = "6k1/1b6/4n3/8/1n4B1/1B3N2/1N6/2b3K1 b - - 0 1"
var CASTLE1FEN = "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"
var CASTLE2FEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
var PERFTFEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

var START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func main() {
	InitAll()

	reader := bufio.NewReader(os.Stdin)

	board := &Board{
		MoveList: &MoveList{},
	}

	board.ParseFen(START_FEN)
	board.GenerateAllMoves()

	board.PrintBoard(os.Stdout)

	_, _ = reader.ReadString('\n')

	// repalce with range
	for moveIndex := 0; moveIndex < board.MoveList.Count; moveIndex++ {
		move := board.MoveList.Moves[moveIndex].Move

		if board.MakeMove(move) == FALSE {
			continue
		}

		fmt.Printf("\nMADE: %s\n", PrintMove(move))
		board.PrintBoard(os.Stdout)

		board.TakeMove()
		fmt.Printf("\nTAKE: %s\n", PrintMove(move))
		board.PrintBoard(os.Stdout)

		_, _ = reader.ReadString('\n')
	}

}

func InitAll() {
	InitSq120To64()
	InitBitMasks()
	InitHashKeys()
	InitFilesRanksBrd()
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
