package engine

import (
	"fmt"
	"os"
	"strings"
)

// FR2SQ converts file and rank coordinates to a square index on a 120-square board.
//
// This function takes file (column) and rank (row) coordinates and calculates the
// corresponding square index on a 120-square board representation. The square index
// is returned as an integer value.
//
// Parameters:
//   - f: The file (column) coordinate, where 'a' corresponds to 0, 'b' to 1, and so on.
//   - r: The rank (row) coordinate, where '1' corresponds to 0, '2' to 1, and so on.
//
// Returns:
//   - An integer representing the square index on a 120-square board.
//
// Example usage:
//
//	squareIndex := FR2SQ(2, 3) // Converts file 'b' and rank '4' to square index.
//
// Note: This function is useful for converting human-readable chess coordinates to
// internal square indices used in the board representation.
func FR2SQ(f, r int) int {
	return (21 + f) + (r * 10)
}

// PrintSq converts a square index to its human-readable algebraic notation.
//
// This function takes a square index on a chessboard (in 120-square format) and
// converts it to its corresponding algebraic notation. The algebraic notation
// represents the square's file (column) and rank (row) on the chessboard, such
// as "a1" for the bottom-left square and "h8" for the top-right square.
// Bewere that our internal representation is reversed, so that a1 is on the top-left
// and h8 is on the bottom right.
//
// Parameters:
//   - sq: The square index on a 120-square board.
//
// Returns:
//   - A string containing the algebraic notation of the square.
//
// Example usage:
//
//	squareNotation := PrintSq(34) // Converts square index 34 to "c4".
//
// Note: This function is useful for displaying chess positions and moves in a
// human-readable format.
func PrintSq(sq int) string {
	file := FilesBrd[sq]
	rank := RanksBrd[sq]

	return fmt.Sprintf("%c%c", ('a' + file), ('1' + rank))
}

// InitFilesRanksBrd initializes the FilesBrd and RanksBrd arrays for square indexing.
//
// This function sets up the `FilesBrd` and `RanksBrd` arrays to enable efficient
// conversion between square indices and their corresponding file (column) and rank
// (row) coordinates. It assigns file and rank values to each square index on a
// 120-square board representation.
//
// The `FilesBrd` array stores the file (column) of each square, and the `RanksBrd`
// array stores the rank (row) of each square. Both arrays are used to convert square
// indices to human-readable algebraic notation and vice versa.
//
// Example usage:
//
//	InitFilesRanksBrd() // Initializes the FilesBrd and RanksBrd arrays.
//	file := FilesBrd[sq] // Retrieves the file (column) of a square.
//	rank := RanksBrd[sq] // Retrieves the rank (row) of a square.
//
// Note: This function should be called once to initialize the arrays before using
// them for square indexing.
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

// UpdateListsMaterial updates the internal lists and material values
// for the given chess position on the board.
//
// This function iterates over all squares on the board, identifies the pieces,
// updates various lists and material counts, and maintains the piece locations.
//
// Parameters:
//   - b: A pointer to the Board structure representing the current chess position.
//
// The following lists and values are updated by this function:
//   - BigPCE: Count of big pieces (rooks and queens) for each color.
//   - MajPCE: Count of major pieces (queens and kings) for each color.
//   - MinPCE: Count of minor pieces (knights and bishops) for each color.
//   - Material: Total material value for each color.
//   - PList: Piece lists containing the positions of each piece type.
//   - PCENum: Count of each piece type on the board.
//   - KingSq: Square positions of the kings for both White and Black.
//   - Pawns: Bitboards representing pawn locations for both colors.
//
// Used to efficiently handles piece type, color, and position tracking
// to support chess engine operations.
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

// Check if sq is onboard
func SqOnBoard(sq int) bool {
	return FilesBrd[sq] != OFFBOARD
}

// Check if sq is offboard
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

// Check if file/rank is valid
func FileRankValid(fileRank int) bool {
	if fileRank >= 0 && fileRank <= 7 {
		return true
	} else {
		return false
	}
}

// Check if piece is valid & empty
func PieceValidEmpty(piece int) bool {
	if piece >= EMPTY && piece <= BK {
		return true
	} else {
		return false
	}
}

// Check if piece is valid
func PieceValid(piece int) bool {
	if piece >= WP && piece <= BK {
		return true
	} else {
		return false
	}
}

func InputWaiting() bool {
	file := os.Stdin
	fi, err := file.Stat()
	if err != nil {
		fmt.Println("file.Stat()", err)
	}
	size := fi.Size()
	if size > 0 {
		return true
	} else {
		return false
	}
}

func ReadInput(s *SearchInfo) {
	var bytes int
	input := make([]byte, 256)

	if InputWaiting() {
		s.Stopped = TRUE
		for bytes < 0 {
			bytes, _ = os.Stdin.Read(input)
		}
		inputStr := string(input[:bytes])
		inputStr = strings.TrimSpace(inputStr)

		if len(inputStr) > 0 && inputStr[:4] == "quit" {
			s.Quit = TRUE
		}
	}
}
