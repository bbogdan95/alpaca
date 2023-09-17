package engine

import "fmt"

type Move struct {
	Move  int
	Score int
}

// most valuable victim - least valuable attacker
var VictimScore = [13]int{0, 100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600}
var MvvLvaScores = [13][13]int{}

func InitMvvLva() {
	for attacker := WP; attacker <= BK; attacker++ {
		for victim := WP; victim <= BK; victim++ {
			MvvLvaScores[victim][attacker] = VictimScore[victim] + 6 - (VictimScore[attacker] / 100)
		}
	}
}

// MoveList is a data structure used to store a list of chess moves.
//
// The MoveList struct represents a list of chess moves, with each move
// being encoded and stored in the 'Moves' array. The 'Count' field keeps
// track of the number of moves in the list. This data structure is
// used for move generation and move ordering within the chess engine.
//
// Fields:
//   - Moves: An array of encoded chess moves.
//   - Count: The number of moves currently in the list.
type MoveList struct {
	Moves [MAXPOSITIONMOVES]Move
	Count int
}

// NewMove constructs a new encoded chess move from its components.
//
// Parameters:
//   - from: The source square of the move (0-119).
//   - to: The destination square of the move (0-119).
//   - captured: The captured piece type (if any) (EMPTY for no capture).
//   - promoted: The promoted piece type (if any) (EMPTY for no promotion).
//   - fileRank: A file and rank combination that can be used for special move flags.
//
// Returns:
//   - An integer representing the encoded chess move.
//
// The NewMove function encodes a chess move by packing the source square,
// destination square, captured piece type, and promoted piece type into a
// single integer. It also allows additional information to be encoded using
// the 'fileRank' parameter, such as flags for en passant, pawn starting move,
// castling, capture, and promotion.
func NewMove(from, to, captured, promoted, fileRank int) int {
	return (((from) | (to << 7) | (captured << 14) | (promoted << 20)) | fileRank)
}

// PrintMoveList prints the contents of a MoveList for debugging or display purposes.
//
// This method displays the number of moves in the MoveList and then iterates
// through each move, printing its index, the human-readable representation of
// the move, and its associated score (if available).
//
// Parameters:
//   - ml: A pointer to the MoveList to be printed.
//
// The printed output includes the move index (1-based), the move in standard
// algebraic notation (e.g., "e2e4" for a pawn move from e2 to e4), and an
// optional score associated with the move. This function is useful for examining
// the list of generated moves during debugging or for displaying move options
// to the user.
func (ml *MoveList) PrintMoveList() {
	fmt.Printf("MoveList: %d\n", ml.Count)

	for i := 0; i < ml.Count; i++ {
		fmt.Printf("Move: %d > %s (score: %d)\n", i+1, PrintMove(ml.Moves[i].Move), ml.Moves[i].Score)
	}

	fmt.Printf("End movelist ----\n\n")
}

// AddQuietMove adds a quiet (non-capturing) chess move to the MoveList.
//
// This method is used to add a quiet move to the MoveList, which typically
// represents pawn or piece moves that do not result in captures. It checks
// whether the source square and destination square of the move are on the
// chessboard. If either of them is off the board, it raises a panic, as
// adding an invalid move is not allowed.
//
// Parameters:
//   - move: The encoded chess move to be added to the MoveList.
//
//
// Note: It's essential to ensure that both source and destination squares of
// the move are on the board before adding the move to the list.
func (ml *MoveList) AddQuietMove(b *Board, move int) {
	if SqOffBoard(GetFrom(move)) || SqOffBoard(GetToSq(move)) {
		panic("cannot add quiet move -- sq offboard")
	}

	ml.Moves[ml.Count].Move = move

	if b.SearchKillers[0][b.Ply] == move {
		ml.Moves[ml.Count].Score = 900000
	} else if b.SearchKillers[1][b.Ply] == move {
		ml.Moves[ml.Count].Score = 800000
	} else {
		ml.Moves[ml.Count].Score = b.SearchHistory[b.Pieces[GetFrom(move)]][GetToSq(move)]
	}

	ml.Count++
}

// AddCaptureMove adds a capturing chess move to the MoveList.
//
// This method is used to add a capturing move to the MoveList, which represents
// moves where one piece captures another. It checks whether the source square
// and destination square of the move are on the chessboard. If either of them
// is off the board, it raises a panic, as adding an invalid move is not allowed.
//
// Parameters:
//   - move: The encoded chess move to be added to the MoveList.
//
// Note: It's essential to ensure that both source and destination squares of
// the move are on the board before adding the move.
func (ml *MoveList) AddCaptureMove(b *Board, move int) {
	if SqOffBoard(GetFrom(move)) || SqOffBoard(GetToSq(move)) {
		panic("cannot add quiet move -- sq offboard (2)")
	}

	ml.Moves[ml.Count].Move = move
	ml.Moves[ml.Count].Score = MvvLvaScores[GetCaptured(move)][b.Pieces[GetFrom(move)]] + 1000000
	ml.Count++
}

// AddEnPassantMove adds an en passant chess move to the MoveList.
//
// This method is used to add an en passant move to the MoveList, which represents
// a special type of capture in chess. It checks whether the source square
// and destination square of the move are on the chessboard. If either of them
// is off the board, it raises a panic, as adding an invalid move is not allowed.
//
// Parameters:
//   - move: The encoded chess move to be added to the MoveList.
//
// Note: It's essential to ensure that both source and destination squares of
// the move are on the board before adding the move.
func (ml *MoveList) AddEnPassantMove(move int) {
	if SqOffBoard(GetFrom(move)) || SqOffBoard(GetToSq(move)) {
		panic("cannot add quiet move -- sq offboard (3)")
	}

	ml.Moves[ml.Count].Move = move
	ml.Moves[ml.Count].Score = 105 + 1000000
	ml.Count++
}

func (ml *MoveList) AddWhitePawnCaptureMove(b *Board, from int, to int, cap int) {
	if !PieceValidEmpty(cap) || SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_7 {
		ml.AddCaptureMove(b, NewMove(from, to, cap, WQ, 0))
		ml.AddCaptureMove(b, NewMove(from, to, cap, WR, 0))
		ml.AddCaptureMove(b, NewMove(from, to, cap, WB, 0))
		ml.AddCaptureMove(b, NewMove(from, to, cap, WN, 0))
	} else {
		ml.AddCaptureMove(b, NewMove(from, to, cap, EMPTY, 0))
	}
}

func (ml *MoveList) AddWhitePawnMove(b *Board, from int, to int) {
	if SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_7 {
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, WQ, 0))
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, WR, 0))
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, WB, 0))
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, WN, 0))
	} else {
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, EMPTY, 0))
	}
}

func (ml *MoveList) AddBlackPawnCaptureMove(b *Board, from int, to int, cap int) {
	if !PieceValidEmpty(cap) || SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_2 {
		ml.AddCaptureMove(b, NewMove(from, to, cap, BQ, 0))
		ml.AddCaptureMove(b, NewMove(from, to, cap, BR, 0))
		ml.AddCaptureMove(b, NewMove(from, to, cap, BB, 0))
		ml.AddCaptureMove(b, NewMove(from, to, cap, BN, 0))
	} else {
		ml.AddCaptureMove(b, NewMove(from, to, cap, EMPTY, 0))
	}
}

func (ml *MoveList) AddBlackPawnMove(b *Board, from int, to int) {
	if SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_2 {
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, BQ, 0))
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, BR, 0))
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, BB, 0))
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, BN, 0))
	} else {
		ml.AddQuietMove(b, NewMove(from, to, EMPTY, EMPTY, 0))
	}
}

// GenerateAllMoves generates all legal chess moves for the current board position
// and populates them in the provided MoveList (ml).
//
// This function is responsible for generating all possible legal moves for the
// current side to move and adding them to the MoveList. It considers various types
// of moves, including pawn moves, en passant captures, castling, sliding piece
// moves, and non-sliding piece moves. The generated moves are suitable for further
// evaluation and searching in the chess engine.
//
// Parameters:
//   - b: A pointer to the Board structure representing the current chess position.
//   - ml: A pointer to the MoveList where the generated moves will be stored.
//
// The generated moves are added to the MoveList, which can then be used for move
// ordering and searching within the chess engine.
//
// Note: This function assumes that the board position is correctly set up, and the
// CheckBoard function is called before generating moves to ensure the board's integrity.
func GenerateAllMoves(b *Board, ml *MoveList) {
	b.CheckBoard()

	ml.Count = 0

	side := b.Side

	// panws
	if side == WHITE {
		for pieceNum := 0; pieceNum < b.PCENum[WP]; pieceNum++ {
			sq := b.PList[WP][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			if b.Pieces[sq+10] == EMPTY {
				ml.AddWhitePawnMove(b, sq, sq+10)
				if RanksBrd[sq] == RANK_2 && b.Pieces[sq+20] == EMPTY {
					ml.AddQuietMove(b, NewMove(sq, sq+20, EMPTY, EMPTY, MoveFlagPawnStart))
				}
			}

			if SqOnBoard(sq+9) && PieceCol[b.Pieces[sq+9]] == BLACK {
				ml.AddWhitePawnCaptureMove(b, sq, sq+9, b.Pieces[sq+9])
			}
			if SqOnBoard(sq+11) && PieceCol[b.Pieces[sq+11]] == BLACK {
				ml.AddWhitePawnCaptureMove(b, sq, sq+11, b.Pieces[sq+11])
			}

			if b.EnPassant != NO_SQ {
				if sq+9 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq+9, EMPTY, EMPTY, MoveFlagEnPassant))
				}
				if sq+11 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq+11, EMPTY, EMPTY, MoveFlagEnPassant))
				}
			}
		}

		if b.CastlePerm&WKCA != 0 {
			if b.Pieces[F1] == EMPTY && b.Pieces[G1] == EMPTY {
				if SqAttacked(E1, BLACK, b) == 0 && SqAttacked(F1, BLACK, b) == 0 {
					ml.AddQuietMove(b, NewMove(E1, G1, EMPTY, EMPTY, MoveFlagCastle))
				}
			}
		}

		if b.CastlePerm&WQCA != 0 {
			if b.Pieces[D1] == EMPTY && b.Pieces[C1] == EMPTY && b.Pieces[B1] == EMPTY {
				if SqAttacked(E1, BLACK, b) == 0 && SqAttacked(D1, BLACK, b) == 0 {
					ml.AddQuietMove(b, NewMove(E1, C1, EMPTY, EMPTY, MoveFlagCastle))
				}
			}
		}
	} else {
		for pieceNum := 0; pieceNum < b.PCENum[BP]; pieceNum++ {
			sq := b.PList[BP][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			if b.Pieces[sq-10] == EMPTY {
				ml.AddBlackPawnMove(b, sq, sq-10)
				if RanksBrd[sq] == RANK_7 && b.Pieces[sq-20] == EMPTY {
					ml.AddQuietMove(b, NewMove(sq, sq-20, EMPTY, EMPTY, MoveFlagPawnStart))
				}
			}

			if SqOnBoard(sq-9) && PieceCol[b.Pieces[sq-9]] == WHITE {
				ml.AddBlackPawnCaptureMove(b, sq, sq-9, b.Pieces[sq-9])
			}
			if SqOnBoard(sq-11) && PieceCol[b.Pieces[sq-11]] == WHITE {
				ml.AddBlackPawnCaptureMove(b, sq, sq-11, b.Pieces[sq-11])
			}

			if b.EnPassant != NO_SQ {
				if sq-9 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq-9, EMPTY, EMPTY, MoveFlagEnPassant))
				}
				if sq-11 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq-11, EMPTY, EMPTY, MoveFlagEnPassant))
				}
			}
		}

		if b.CastlePerm&BKCA != 0 {
			if b.Pieces[F8] == EMPTY && b.Pieces[G8] == EMPTY {
				if SqAttacked(E8, WHITE, b) == 0 && SqAttacked(F8, WHITE, b) == 0 {
					ml.AddQuietMove(b, NewMove(E8, G8, EMPTY, EMPTY, MoveFlagCastle))
				}
			}
		}

		if b.CastlePerm&BQCA != 0 {
			if b.Pieces[D8] == EMPTY && b.Pieces[C8] == EMPTY && b.Pieces[B8] == EMPTY {
				if SqAttacked(E8, WHITE, b) == 0 && SqAttacked(D8, WHITE, b) == 0 {
					ml.AddQuietMove(b, NewMove(E8, C8, EMPTY, EMPTY, MoveFlagCastle))
				}
			}
		}
	}

	// slide pieces
	pieceIndex := LoopSlideIndex[side]
	piece := LoopSlidePieces[pieceIndex]
	pieceIndex++
	for piece != 0 {
		if !PieceValid(piece) {
			panic("piece not valid")
		}

		for pieceNum := 0; pieceNum < b.PCENum[piece]; pieceNum++ {
			sq := b.PList[piece][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			for i := 0; i < NumDir[piece]; i++ {
				dir := PieceDir[piece][i]
				tSq := sq + dir

				for SqOnBoard(tSq) {
					if b.Pieces[tSq] != EMPTY {
						if PieceCol[b.Pieces[tSq]] == side^1 {
							ml.AddCaptureMove(b, NewMove(sq, tSq, b.Pieces[tSq], EMPTY, 0))
						}
						break
					}

					ml.AddQuietMove(b, NewMove(sq, tSq, EMPTY, EMPTY, 0))
					tSq += dir
				}

			}
		}

		piece = LoopSlidePieces[pieceIndex]
		pieceIndex++
	}

	// non slide pieces
	pieceIndex = LoopNonSlideIndex[side]
	piece = LoopNonSlidePieces[pieceIndex]
	pieceIndex++
	for piece != 0 {
		if !PieceValid(piece) {
			panic("piece not valid")
		}

		for pieceNum := 0; pieceNum < b.PCENum[piece]; pieceNum++ {
			sq := b.PList[piece][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			for i := 0; i < NumDir[piece]; i++ {
				dir := PieceDir[piece][i]
				tSq := sq + dir

				if SqOffBoard(tSq) {
					continue
				}

				if b.Pieces[tSq] != EMPTY {
					if PieceCol[b.Pieces[tSq]] == side^1 {
						ml.AddCaptureMove(b, NewMove(sq, tSq, b.Pieces[tSq], EMPTY, 0))
					}
					continue
				}

				ml.AddQuietMove(b, NewMove(sq, tSq, EMPTY, EMPTY, 0))
			}
		}

		piece = LoopNonSlidePieces[pieceIndex]
		pieceIndex++
	}
}

func GenerateAllCaptures(b *Board, ml *MoveList) {
	b.CheckBoard()

	ml.Count = 0

	side := b.Side

	// panws
	if side == WHITE {
		for pieceNum := 0; pieceNum < b.PCENum[WP]; pieceNum++ {
			sq := b.PList[WP][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			if SqOnBoard(sq+9) && PieceCol[b.Pieces[sq+9]] == BLACK {
				ml.AddWhitePawnCaptureMove(b, sq, sq+9, b.Pieces[sq+9])
			}
			if SqOnBoard(sq+11) && PieceCol[b.Pieces[sq+11]] == BLACK {
				ml.AddWhitePawnCaptureMove(b, sq, sq+11, b.Pieces[sq+11])
			}

			if b.EnPassant != NO_SQ {
				if sq+9 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq+9, EMPTY, EMPTY, MoveFlagEnPassant))
				}
				if sq+11 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq+11, EMPTY, EMPTY, MoveFlagEnPassant))
				}
			}
		}
	} else {
		for pieceNum := 0; pieceNum < b.PCENum[BP]; pieceNum++ {
			sq := b.PList[BP][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			if SqOnBoard(sq-9) && PieceCol[b.Pieces[sq-9]] == WHITE {
				ml.AddBlackPawnCaptureMove(b, sq, sq-9, b.Pieces[sq-9])
			}
			if SqOnBoard(sq-11) && PieceCol[b.Pieces[sq-11]] == WHITE {
				ml.AddBlackPawnCaptureMove(b, sq, sq-11, b.Pieces[sq-11])
			}

			if b.EnPassant != NO_SQ {
				if sq-9 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq-9, EMPTY, EMPTY, MoveFlagEnPassant))
				}
				if sq-11 == b.EnPassant {
					ml.AddEnPassantMove(NewMove(sq, sq-11, EMPTY, EMPTY, MoveFlagEnPassant))
				}
			}
		}
	}

	// slide pieces
	pieceIndex := LoopSlideIndex[side]
	piece := LoopSlidePieces[pieceIndex]
	pieceIndex++
	for piece != 0 {
		if !PieceValid(piece) {
			panic("piece not valid")
		}

		for pieceNum := 0; pieceNum < b.PCENum[piece]; pieceNum++ {
			sq := b.PList[piece][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			for i := 0; i < NumDir[piece]; i++ {
				dir := PieceDir[piece][i]
				tSq := sq + dir

				for SqOnBoard(tSq) {
					if b.Pieces[tSq] != EMPTY {
						if PieceCol[b.Pieces[tSq]] == side^1 {
							ml.AddCaptureMove(b, NewMove(sq, tSq, b.Pieces[tSq], EMPTY, 0))
						}
						break
					}
					tSq += dir
				}

			}
		}

		piece = LoopSlidePieces[pieceIndex]
		pieceIndex++
	}

	// non slide pieces
	pieceIndex = LoopNonSlideIndex[side]
	piece = LoopNonSlidePieces[pieceIndex]
	pieceIndex++
	for piece != 0 {
		if !PieceValid(piece) {
			panic("piece not valid")
		}

		for pieceNum := 0; pieceNum < b.PCENum[piece]; pieceNum++ {
			sq := b.PList[piece][pieceNum]
			if SqOffBoard(sq) {
				panic("sq offboard")
			}

			for i := 0; i < NumDir[piece]; i++ {
				dir := PieceDir[piece][i]
				tSq := sq + dir

				if SqOffBoard(tSq) {
					continue
				}

				if b.Pieces[tSq] != EMPTY {
					if PieceCol[b.Pieces[tSq]] == side^1 {
						ml.AddCaptureMove(b, NewMove(sq, tSq, b.Pieces[tSq], EMPTY, 0))
					}
					continue
				}
			}
		}

		piece = LoopNonSlidePieces[pieceIndex]
		pieceIndex++
	}
}

func MoveExists(b *Board, move int) int {
	var ml MoveList
	GenerateAllMoves(b, &ml)
	for i := 0; i < ml.Count; i++ {
		res, err := b.MakeMove(move)
		if err != nil {
			panic(err)
		}

		if res == 0 {
			continue
		}

		b.TakeMove()
		if ml.Moves[i].Move == move {
			return TRUE
		}
	}

	return FALSE
}
