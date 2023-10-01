package engine

import (
	"errors"
	"fmt"
)

/*
A chess move is represented as a 32-bit integer (int32).

Bits layout:
  0000 0000 0000 0000 0000 0111 1111 -> Bits 0-6: 'from' square index.
  0000 0000 0000 0011 1111 1000 0000 -> Bits 7-13: 'to' square index.
  0000 0000 0011 1100 0000 0000 0000 -> Bits 14-17: Capture move flag with captured piece index.
  0000 0000 0100 0000 0000 0000 0000 -> Bit 18: En passant capture move flag.
  0000 0111 1000 0000 0000 0000 0000 -> Bits 19-22: Piece promotion flag with promoted piece index.
  0000 1000 0000 0000 0000 0000 0000 -> Bit 23: Castling move flag.

To create a new move, use the NewMove function, e.g., NewMove(from, to, EMPTY, EMPTY, 0).
*/

// Move flag En Passant
var MoveFlagEnPassant = 0x40000

// Move flag Pawn Start
var MoveFlagPawnStart = 0x80000

// Move flag Castle
var MoveFlagCastle = 0x1000000

// Move flag Capture
var MoveFlagCapture = 0x7C000

// Move flag Promoted
var MoveFlagPromotion = 0xF00000

var NOMOVE = 0

// Extract the 'from' part of a move
func GetFrom(move int) int {
	return move & 0x7F
}

// Extract the 'to' part of a move
func GetToSq(move int) int {
	return (move >> 7) & 0x7F
}

// Extract the captured
func GetCaptured(move int) int {
	return (move >> 14) & 0xF
}
func GetPromoted(move int) int {
	return (move >> 20) & 0xF
}

// PrintMove converts an encoded chess move to human-readable algebraic notation.
//
// This function takes an encoded chess move (integer) and converts it into a human-readable
// algebraic notation that represents the move. The algebraic notation typically includes
// the source square, destination square, and, if applicable, the piece promotion information.
//
// Parameters:
//   - move: An encoded chess move, as an integer.
//
// Returns:
//   - A string containing the human-readable algebraic notation of the move.
//
// Example usage:
//
//	moveStr := PrintMove(encodedMove) // Converts the encoded move to algebraic notation.
//
// Note: This function handles both standard moves and moves with piece promotions, such as
// pawn promotions to queen, rook, bishop, or knight.
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

// ParseMove parses a chess move from the given algebraic notation string
func ParseMove(move string, b *Board) (int, error) {
	if len(move) < 4 {
		return NOMOVE, nil
	}
	if move[1] > '8' || move[1] < '1' {
		return NOMOVE, nil
	}
	if move[3] > '8' || move[3] < '1' {
		return NOMOVE, nil
	}
	if move[0] > 'h' || move[0] < 'a' {
		return NOMOVE, nil
	}
	if move[2] > 'h' || move[2] < 'a' {
		return NOMOVE, nil
	}

	from := FR2SQ(int(move[0]-'a'), int(move[1]-'1'))
	to := FR2SQ(int(move[2]-'a'), int(move[3]-'1'))

	if SqOffBoard(from) || SqOffBoard(to) {
		return 0, errors.New("sq is offboard")
	}

	var ml MoveList
	GenerateAllMoves(b, &ml)

	promotedPiece := EMPTY

	for i := 0; i < ml.Count; i++ {
		m := ml.Moves[i].Move
		if GetFrom(m) == from && GetToSq(m) == to {
			promotedPiece = GetPromoted(m)
			if promotedPiece != EMPTY {
				if PieceRookQueen[promotedPiece] != 0 && PieceBishopQueen[promotedPiece] == 0 && move[4] == 'r' {
					return m, nil
				} else if PieceRookQueen[promotedPiece] == 0 && PieceBishopQueen[promotedPiece] != 0 && move[4] == 'b' {
					return m, nil
				} else if PieceRookQueen[promotedPiece] != 0 && PieceBishopQueen[promotedPiece] != 0 && move[4] == 'q' {
					return m, nil
				} else if PieceKnight[promotedPiece] != 0 && move[4] == 'n' {
					return m, nil
				}
				continue
			}

			return m, nil
		}
	}

	return NOMOVE, nil
}
