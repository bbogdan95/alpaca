package engine

import "math"

var PawnIsolated = -10
var PawnPassed = [8]int{0, 5, 10, 20, 35, 60, 100, 200}
var RookOpenFile = 10
var RookSemiOpenFile = 5
var QueenOpenFile = 5
var QueenSemiOpenFile = 3
var BishopPair = 30

// used to switch between KingO for openings and KingE for endgames
var EndGameMaterial = 1*PieceVal[WR] + 2*PieceVal[WN] + 2*PieceVal[WP]

var PawnTable = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 0, -10, -10, 0, 10, 10,
	5, 0, 0, 5, 5, 0, 0, 5,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 5, 5, 10, 10, 5, 5, 5,
	10, 10, 10, 20, 20, 10, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KnightTable = [64]int{
	0, -10, 0, 0, 0, 0, -10, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 5, 0,
	5, 10, 15, 20, 20, 15, 10, 5,
	5, 10, 10, 20, 20, 10, 10, 5,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var BishopTable = [64]int{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var RookTable = [64]int{
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	25, 25, 25, 25, 25, 25, 25, 25,
	0, 0, 5, 10, 10, 5, 0, 0,
}

// in the end game, the king is encouraged to move to the
// center of the board, and is penalised from moving to the corners or the sides
var KingE = [64]int{
	-50, -10, 0, 0, 0, 0, -10, -50,
	-10, 0, 10, 10, 10, 10, 0, -10,
	0, 10, 20, 20, 20, 20, 10, 0,
	0, 10, 20, 40, 40, 20, 10, 0,
	0, 10, 20, 40, 40, 20, 10, 0,
	0, 10, 20, 20, 20, 20, 10, 0,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-50, -10, 0, 0, 0, 0, -10, -50,
}

// in the opening the king is encouraged to castle
// and is heavely penalised to advanced forward
var KingO = [64]int{
	0, 5, 5, -10, -10, 0, 10, 5,
	-30, -30, -30, -30, -30, -30, -30, -30,
	-50, -50, -50, -50, -50, -50, -50, -50,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
}

var Mirror64 = [64]int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

func EvalPosition(b *Board) int {
	score := b.Material[WHITE] - b.Material[BLACK]

	if b.PCENum[WP] == 0 && b.PCENum[BP] == 0 && MaterialDraw(b) == TRUE {
		return 0
	}

	piece := WP
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score += PawnTable[SQ64[sq]]

		if IsolatedMask[SQ64[sq]]&b.Pawns[WHITE] == 0 {
			score += PawnIsolated
		}

		if WhitePassedMask[SQ64[sq]]&b.Pawns[BLACK] == 0 {
			score += PawnPassed[RanksBrd[sq]]
		}
	}

	piece = BP
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score -= PawnTable[Mirror64[SQ64[sq]]]

		if IsolatedMask[SQ64[sq]]&b.Pawns[BLACK] == 0 {
			score -= PawnIsolated
		}

		if BlackPassedMask[SQ64[sq]]&b.Pawns[WHITE] == 0 {
			score -= PawnPassed[7-RanksBrd[sq]]
		}
	}

	piece = WN
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score += KnightTable[SQ64[sq]]
	}

	piece = BN
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score -= KnightTable[Mirror64[SQ64[sq]]]
	}

	piece = WB
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score += BishopTable[SQ64[sq]]
	}

	piece = BB
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score -= BishopTable[Mirror64[SQ64[sq]]]
	}

	piece = WR
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score += RookTable[SQ64[sq]]

		if b.Pawns[BOTH]&FileBBMask[FilesBrd[sq]] == 0 {
			score += RookOpenFile
		} else if b.Pawns[WHITE]&FileBBMask[FilesBrd[sq]] == 0 {
			score += RookSemiOpenFile
		}
	}

	piece = BR
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		score -= RookTable[Mirror64[SQ64[sq]]]

		if b.Pawns[BOTH]&FileBBMask[FilesBrd[sq]] == 0 {
			score -= RookOpenFile
		} else if b.Pawns[BLACK]&FileBBMask[FilesBrd[sq]] == 0 {
			score -= RookSemiOpenFile
		}
	}

	piece = WQ
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		if b.Pawns[BOTH]&FileBBMask[FilesBrd[sq]] == 0 {
			score += QueenOpenFile
		} else if b.Pawns[WHITE]&FileBBMask[FilesBrd[sq]] == 0 {
			score += QueenSemiOpenFile
		}
	}

	piece = BQ
	for i := 0; i < b.PCENum[piece]; i++ {
		sq := b.PList[piece][i]
		if b.Pawns[BOTH]&FileBBMask[FilesBrd[sq]] == 0 {
			score -= QueenOpenFile
		} else if b.Pawns[BLACK]&FileBBMask[FilesBrd[sq]] == 0 {
			score -= QueenSemiOpenFile
		}
	}

	piece = WK
	sq := b.PList[piece][0]
	if b.Material[BLACK] <= EndGameMaterial {
		score += KingE[SQ64[sq]]
	} else {
		score += KingO[SQ64[sq]]
	}

	piece = BK
	sq = b.PList[piece][0]
	if b.Material[WHITE] <= EndGameMaterial {
		score -= KingE[Mirror64[SQ64[sq]]]
	} else {
		score -= KingO[Mirror64[SQ64[sq]]]
	}

	if b.PCENum[WB] >= 2 {
		score += BishopPair
	}
	if b.PCENum[BB] >= 2 {
		score -= BishopPair
	}

	if b.Side == WHITE {
		return score
	} else {
		return -score
	}
}

// sjeng 11.2
func MaterialDraw(b *Board) int {
	if b.PCENum[WR] == 0 && b.PCENum[BR] == 0 && b.PCENum[WQ] == 0 && b.PCENum[BQ] == 0 {
		if b.PCENum[BB] == 0 && b.PCENum[WB] == 0 {
			if b.PCENum[WN] < 3 && b.PCENum[BN] < 3 {
				return 1
			}
		} else if b.PCENum[WN] == 0 && b.PCENum[BN] == 0 {
			if math.Abs(float64(b.PCENum[WB]-b.PCENum[BB])) < 2 {
				return 1
			}
		} else if (b.PCENum[WN] < 3 && b.PCENum[WB] == 0) || (b.PCENum[WB] == 1 && b.PCENum[WN] == 0) {
			if (b.PCENum[BN] < 3 && b.PCENum[BB] == 0) || (b.PCENum[BB] == 1 && b.PCENum[BN] == 0) {
				return 1
			}
		}
	} else if b.PCENum[WQ] == 0 && b.PCENum[BQ] == 0 {
		if b.PCENum[WR] == 1 && b.PCENum[BR] == 1 {
			if (b.PCENum[WN]+b.PCENum[WB]) < 2 && (b.PCENum[BN]+b.PCENum[BB]) < 2 {
				return 1
			}
		} else if b.PCENum[WR] == 1 && b.PCENum[BR] == 0 {
			if (b.PCENum[WN]+b.PCENum[WB] == 0) && (((b.PCENum[BN] + b.PCENum[BB]) == 1) || ((b.PCENum[BN] + b.PCENum[BB]) == 2)) {
				return TRUE
			}
		} else if b.PCENum[BR] == 1 && b.PCENum[WR] == 0 {
			if (b.PCENum[BN]+b.PCENum[BB] == 0) && (((b.PCENum[WN] + b.PCENum[WB]) == 1) || ((b.PCENum[WN] + b.PCENum[WB]) == 2)) {
				return TRUE
			}
		}
	}

	return 0
}
