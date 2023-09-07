package main

import "fmt"

func FR2SQ(f, r int) int {
	return (21 + f) + (r * 10)
}

// Move flag En Passant
var MFLAGEP = 0x40000

// Move flag Pawn Start
var MFLAGPS = 0x80000

// Move flag Castle
var MFLAGCA = 0x100000

// Move flag Capture
var MFLAGCAP = 0x7c0000

// Move flag Promoted
var MFLAGPROM = 0xF0000

func GetFrom(move int) int {
	return move & 0x7F
}
func GetToSq(move int) int {
	return (move >> 7) & 0x7F
}
func GetCaptured(move int) int {
	return (move >> 14) & 0xF
}
func GetPromoted(move int) int {
	return (move >> 20) & 0xF
}

// Takes a square in the 120 board square representation
// and returns its string notation. (eq. sq 21 is a1)
func PrintSq(sq int) string {
	file := FilesBrd[sq]
	rank := RanksBrd[sq]

	return fmt.Sprintf("%c%c", ('a' + file), ('1' + rank))
}

// Prints algebraic notation string move
func PrintMove(move int) string {
	fileFrom := FilesBrd[GetFrom(move)]
	rankFrom := RanksBrd[GetFrom(move)]

	fileTo := FilesBrd[GetToSq(move)]
	rankTo := RanksBrd[GetToSq(move)]

	promoted := GetPromoted(move)

	if promoted != 0 {
		pieceRune := 'q'
		if PieceKnight[promoted] == TRUE {
			pieceRune = 'n'
		} else if PieceRookQueen[promoted] == TRUE && PieceBishopQueen[promoted] == FALSE {
			pieceRune = 'r'
		} else if PieceRookQueen[promoted] == FALSE && PieceBishopQueen[promoted] == TRUE {
			pieceRune = 'b'
		}

		return fmt.Sprintf("%c%c%c%c%c", ('a' + fileFrom), ('1' + rankFrom), ('a' + fileTo), ('1' + rankTo), pieceRune)
	} else {
		return fmt.Sprintf("%c%c%c%c", ('a' + fileFrom), ('1' + rankFrom), ('a' + fileTo), ('1' + rankTo))
	}
}

func InitFilesRanksBrd() {
	for i := 0; i < BRD_SQ_NUM; i++ {
		FilesBrd[i] = OFFBOARD
		RanksBrd[i] = OFFBOARD
	}

	for rank := RANK_1; rank <= RANK_8; rank++ {
		for file := FILE_A; file <= FILE_H; file++ {
			sq := FR2SQ(file, rank)
			FilesBrd[sq] = file
			RanksBrd[sq] = rank
		}
	}
}

func UpdateListsMaterial(b *Board) {
	for i := 0; i < BRD_SQ_NUM; i++ {
		sq := i
		piece := b.Pieces[sq]
		if piece != OFFBOARD && piece != EMPTY {
			color := PieceCol[piece]

			if PieceBig[piece] == TRUE {
				b.BigPCE[color]++
			}
			if PieceMaj[piece] == TRUE {
				b.MajPCE[color]++
			}
			if PieceMin[piece] == TRUE {
				b.MinPCE[color]++
			}

			b.Material[color] += PieceVal[piece]

			b.PList[piece][b.PCENum[piece]] = sq
			b.PCENum[piece]++

			if piece == WK {
				b.KingSq[WHITE] = sq
			}
			if piece == BK {
				b.KingSq[BLACK] = sq
			}

			if piece == WP {
				SetBit(&b.Pawns[WHITE], SQ64[sq])
				SetBit(&b.Pawns[BOTH], SQ64[sq])
			} else if piece == BP {
				SetBit(&b.Pawns[BLACK], SQ64[sq])
				SetBit(&b.Pawns[BOTH], SQ64[sq])
			}
		}
	}
}

// Check is sq is onboard
func SqOnBoard(sq int) bool {
	return FilesBrd[sq] != OFFBOARD
}

func SqOffBoard(sq int) bool {
	return FilesBrd[sq] == OFFBOARD
}

// Check if side is valid
func SideValid(side int) bool {
	if side == WHITE || side == BLACK {
		return true
	} else {
		return false
	}
}

func FileRankValid(fileRank int) bool {
	if fileRank >= 0 && fileRank <= 7 {
		return true
	} else {
		return false
	}
}

func PieceValidEmpty(piece int) bool {
	if piece >= EMPTY && piece <= BK {
		return true
	} else {
		return false
	}
}

func PieceValid(piece int) bool {
	if piece >= WP && piece <= BK {
		return true
	} else {
		return false
	}
}
