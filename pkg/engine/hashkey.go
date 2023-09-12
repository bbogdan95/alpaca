package engine

import "math/rand"

/*
The Zobrist key is a concept used in computer chess programming to efficiently represent and manipulate chess positions.
It's named after the computer scientist and chess programmer Alberto Zobrist, who introduced the idea.
The Zobrist key is essentially a random 64-bit integer associated with each possible combination of pieces on each square of the chessboard,
as well as other factors that can affect the state of the game, such as the side to move, castling rights, and the en passant square.
These random keys are precomputed and stored in tables, typically referred to as Zobrist tables.
*/

// PieceKeys is an array used to store random 64-bit integers for each piece type (including empty squares) on each square of the chessboard.
// In the InitHashKeys function, it generates random keys for all 13 piece types (including an empty square) on all 120 squares (including off-board squares).
// These random keys are used to calculate the hash key for each position, known as Zobrist keys.
// The Zobrist keys are XORed together to create the hash key for a given position.
var PieceKeys [13][120]uint64

// SideKey is a random key used to represent which side (White or Black) is to move in the current position.
// In the InitHashKeys function, it generates a random key for the side to move.
// This key is XORed with the hash key when it's White's turn to move, effectively flipping the bit to indicate the change of side.
var SideKey uint64

// CastleKeys is an array used to represent the castling rights in a position.
// In the InitHashKeys function, it generates random keys for all possible castling rights combinations (16 in total).
// These keys are XORed with the hash key to account for the castling rights in a position.
var CastleKeys [16]uint64

// InitHashKeys initializes the Zobrist hash keys used for chess board positions.
//
// This function populates the PieceKeys array with random 64-bit integers for each piece type (including
// empty squares) on each square of the 120-square board representation. It also generates a SideKey,
// representing the side to move, and CastleKeys, representing the castling rights.
//
// Note: It is essential to call this function at the start of the chess engine to initialize the Zobrist
// keys for later use.
//
// Example usage:
//   InitHashKeys() // Initializes Zobrist keys for position hashing.
func InitHashKeys() {
	for i := 0; i < 13; i++ {
		for j := 0; j < 120; j++ {
			PieceKeys[i][j] = rand.Uint64()
		}
	}

	SideKey = rand.Uint64()
	for i := 0; i < 16; i++ {
		CastleKeys[i] = rand.Uint64()
	}
}

// GeneratePosKey calculates and returns the position key (Zobrist key) for the current chess board state.
//
// The position key is a fundamental component of chess engines, serving as a unique identifier for a specific
// chess position. It is computed using Zobrist hashing, a technique that combines random keys for various
// aspects of the board state. These aspects include the piece placements, the side to move, the en passant
// square (if any), and the castling rights.
//
// Parameters:
//   - b: A pointer to the current chess board (type *Board) representing the board state.
//
// Returns:
//   - The position key (Zobrist key) as a 64-bit unsigned integer (uint64).
//
// Details:
//   - The position key is calculated by XOR-ing together the following components:
//     1. Piece Keys: For each square on the board where there is a piece (excluding off-board squares),
//        the corresponding PieceKeys value is XOR-ed into the finalKey.
//     2. Side Key: If it is White's turn to move, the SideKey is XOR-ed into the finalKey.
//     3. En Passant Key: If there is a valid en passant square, the PieceKeys value for an empty square at
//        that location is XOR-ed into the finalKey.
//     4. Castle Key: The CastleKeys value corresponding to the current castling rights is XOR-ed into the finalKey.
//
// Example usage:
//   key := GeneratePosKey(board) // Calculates the position key for the current board state.
func GeneratePosKey(b *Board) uint64 {
	var finalKey uint64 = 0
	piece := EMPTY

	for sq := 0; sq < BRD_SQ_NUM; sq++ {
		piece = b.Pieces[sq]
		if piece != NO_SQ && piece != EMPTY && piece != OFFBOARD {
			finalKey ^= PieceKeys[piece][sq]
		}
	}

	if b.Side == WHITE {
		finalKey ^= SideKey
	}

	if b.EnPassant != NO_SQ {
		finalKey ^= PieceKeys[EMPTY][b.EnPassant]
	}

	finalKey ^= CastleKeys[b.CastlePerm]

	return finalKey
}
