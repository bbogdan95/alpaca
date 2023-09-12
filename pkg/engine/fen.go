package engine

// ParseFen parses a FEN (Forsyth-Edwards Notation) string and sets up the board
// according to the specified position.
//
// This method takes a FEN string and configures the internal state of the Board
// structure to match the specified position. FEN is a standard notation used to
// describe a chess position, including piece placement, castling rights, en passant
// target square, and other game state information.
//
// Parameters:
//   - fen: A string containing the FEN representation of the chess position.
//
// The method iterates through the FEN string, interpreting its components, and sets
// up the board accordingly. It also computes the position key (Zobrist key) and
// updates various lists and material counts based on the position.
//
// Example usage:
//   var board Board
//   board.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
//   // The board is now set up with the specified position.
//
// Note: The input FEN string should be correctly formatted and represent a valid
// chess position.
func (b *Board) ParseFen(fen string) {
	rank := RANK_8
	file := FILE_A
	piece := 0
	count := 0
	sq64 := 0
	sq120 := 0
	fenPos := 0

	b.ResetBoard()

	for i, p := range fen {
		if rank < RANK_1 {
			break
		}
		fenPos = i
		count = 1
		switch p {
		case 'p':
			piece = BP
		case 'r':
			piece = BR
		case 'n':
			piece = BN
		case 'b':
			piece = BB
		case 'k':
			piece = BK
		case 'q':
			piece = BQ
		case 'P':
			piece = WP
		case 'R':
			piece = WR
		case 'N':
			piece = WN
		case 'B':
			piece = WB
		case 'K':
			piece = WK
		case 'Q':
			piece = WQ
		case '1', '2', '3', '4', '5', '6', '7', '8':
			piece = EMPTY
			count = int(p - '0')
		case '/', ' ':
			rank--
			file = FILE_A
			fenPos++
			continue
		default:
			panic("FEN error")
		}

		for i := 0; i < count; i++ {
			sq64 = rank*8 + file
			sq120 = SQ120[sq64]
			if piece != EMPTY {
				b.Pieces[sq120] = piece
			}
			file++
		}
	}

	if fen[fenPos] == 'w' {
		b.Side = WHITE
	} else {
		b.Side = BLACK
	}

	fenPos += 2

	for i := 0; i < 4; i++ {
		if fen[fenPos] == ' ' {
			break
		}
		switch fen[fenPos] {
		case 'K':
			b.CastlePerm |= WKCA
		case 'Q':
			b.CastlePerm |= WQCA
		case 'k':
			b.CastlePerm |= BKCA
		case 'q':
			b.CastlePerm |= BQCA
		}

		fenPos++
	}
	fenPos++

	if fen[fenPos] != '-' {
		file = int(fen[fenPos] - 'a')
		fenPos++
		rank = int(fen[fenPos] - '1')

		b.EnPassant = FR2SQ(file, rank)
	}

	b.PosKey = GeneratePosKey(b)
	UpdateListsMaterial(b)
}
