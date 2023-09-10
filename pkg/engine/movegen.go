package engine

import "fmt"

type Move struct {
	Move  int
	Score int
}

type MoveList struct {
	Moves [MAXPOSITIONMOVES]Move
	// number of moves
	Count int
}

func NewMove(from, to, captured, promoted, fileRank int) int {
	return (((from) | (to << 7) | (captured << 14) | (promoted << 20)) | fileRank)
}

func (ml *MoveList) PrintMoveList() {
	fmt.Printf("MoveList: %d\n", ml.Count)

	for i := 0; i < ml.Count; i++ {
		fmt.Printf("Move: %d > %s (score: %d)\n", i+1, PrintMove(ml.Moves[i].Move), ml.Moves[i].Score)
	}

	fmt.Printf("End movelist ----\n\n")
}

func (ml *MoveList) AddQuietMove(move int) {
	if SqOffBoard(GetFrom(move)) || SqOffBoard(GetToSq(move)) {
		panic("cannot add quiet move -- sq offboard")
	}

	ml.Moves[ml.Count].Move = move
	ml.Moves[ml.Count].Score = 0
	ml.Count++
}

func (ml *MoveList) AddCaptureMove(move int) {
	if SqOffBoard(GetFrom(move)) || SqOffBoard(GetToSq(move)) {
		panic("cannot add quiet move -- sq offboard (2)")
	}

	ml.Moves[ml.Count].Move = move
	ml.Moves[ml.Count].Score = 0
	ml.Count++
}

func (ml *MoveList) AddEnPassantMove(move int) {
	if SqOffBoard(GetFrom(move)) || SqOffBoard(GetToSq(move)) {
		panic("cannot add quiet move -- sq offboard (3)")
	}

	ml.Moves[ml.Count].Move = move
	ml.Moves[ml.Count].Score = 0
	ml.Count++
}

func (ml *MoveList) AddWhitePawnCaptureMove(from int, to int, cap int) {
	if !PieceValidEmpty(cap) || SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_7 {
		ml.AddCaptureMove(NewMove(from, to, cap, WQ, 0))
		ml.AddCaptureMove(NewMove(from, to, cap, WR, 0))
		ml.AddCaptureMove(NewMove(from, to, cap, WB, 0))
		ml.AddCaptureMove(NewMove(from, to, cap, WN, 0))
	} else {
		ml.AddCaptureMove(NewMove(from, to, cap, EMPTY, 0))
	}
}

func (ml *MoveList) AddWhitePawnMove(from int, to int) {
	if SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_7 {
		ml.AddQuietMove(NewMove(from, to, EMPTY, WQ, 0))
		ml.AddQuietMove(NewMove(from, to, EMPTY, WR, 0))
		ml.AddQuietMove(NewMove(from, to, EMPTY, WB, 0))
		ml.AddQuietMove(NewMove(from, to, EMPTY, WN, 0))
	} else {
		ml.AddQuietMove(NewMove(from, to, EMPTY, EMPTY, 0))
	}
}

func (ml *MoveList) AddBlackPawnCaptureMove(from int, to int, cap int) {
	if !PieceValidEmpty(cap) || SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_2 {
		ml.AddCaptureMove(NewMove(from, to, cap, BQ, 0))
		ml.AddCaptureMove(NewMove(from, to, cap, BR, 0))
		ml.AddCaptureMove(NewMove(from, to, cap, BB, 0))
		ml.AddCaptureMove(NewMove(from, to, cap, BN, 0))
	} else {
		ml.AddCaptureMove(NewMove(from, to, cap, EMPTY, 0))
	}
}

func (ml *MoveList) AddBlackPawnMove(from int, to int) {
	if SqOffBoard(from) || SqOffBoard(to) {
		panic("assert failed (1)")
	}

	if RanksBrd[from] == RANK_2 {
		ml.AddQuietMove(NewMove(from, to, EMPTY, BQ, 0))
		ml.AddQuietMove(NewMove(from, to, EMPTY, BR, 0))
		ml.AddQuietMove(NewMove(from, to, EMPTY, BB, 0))
		ml.AddQuietMove(NewMove(from, to, EMPTY, BN, 0))
	} else {
		ml.AddQuietMove(NewMove(from, to, EMPTY, EMPTY, 0))
	}
}

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
				ml.AddWhitePawnMove(sq, sq+10)
				if RanksBrd[sq] == RANK_2 && b.Pieces[sq+20] == EMPTY {
					ml.AddQuietMove(NewMove(sq, sq+20, EMPTY, EMPTY, MoveFlagPawnStart))
				}
			}

			if SqOnBoard(sq+9) && PieceCol[b.Pieces[sq+9]] == BLACK {
				ml.AddWhitePawnCaptureMove(sq, sq+9, b.Pieces[sq+9])
			}
			if SqOnBoard(sq+11) && PieceCol[b.Pieces[sq+11]] == BLACK {
				ml.AddWhitePawnCaptureMove(sq, sq+11, b.Pieces[sq+11])
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
					ml.AddQuietMove(NewMove(E1, G1, EMPTY, EMPTY, MoveFlagCastle))
				}
			}
		}

		if b.CastlePerm&WQCA != 0 {
			if b.Pieces[D1] == EMPTY && b.Pieces[C1] == EMPTY && b.Pieces[B1] == EMPTY {
				if SqAttacked(E1, BLACK, b) == 0 && SqAttacked(D1, BLACK, b) == 0 {
					ml.AddQuietMove(NewMove(E1, C1, EMPTY, EMPTY, MoveFlagCastle))
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
				ml.AddBlackPawnMove(sq, sq-10)
				if RanksBrd[sq] == RANK_7 && b.Pieces[sq-20] == EMPTY {
					ml.AddQuietMove(NewMove(sq, sq-20, EMPTY, EMPTY, MoveFlagPawnStart))
				}
			}

			if SqOnBoard(sq-9) && PieceCol[b.Pieces[sq-9]] == WHITE {
				ml.AddBlackPawnCaptureMove(sq, sq-9, b.Pieces[sq-9])
			}
			if SqOnBoard(sq-11) && PieceCol[b.Pieces[sq-11]] == WHITE {
				ml.AddBlackPawnCaptureMove(sq, sq-11, b.Pieces[sq-11])
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
					ml.AddQuietMove(NewMove(E8, G8, EMPTY, EMPTY, MoveFlagCastle))
				}
			}
		}

		if b.CastlePerm&BQCA != 0 {
			if b.Pieces[D8] == EMPTY && b.Pieces[C8] == EMPTY && b.Pieces[B8] == EMPTY {
				if SqAttacked(E8, WHITE, b) == 0 && SqAttacked(D8, WHITE, b) == 0 {
					ml.AddQuietMove(NewMove(E8, C8, EMPTY, EMPTY, MoveFlagCastle))
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
							ml.AddCaptureMove(NewMove(sq, tSq, b.Pieces[tSq], EMPTY, 0))
						}
						break
					}

					ml.AddQuietMove(NewMove(sq, tSq, EMPTY, EMPTY, 0))
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
						ml.AddCaptureMove(NewMove(sq, tSq, b.Pieces[tSq], EMPTY, 0))
					}
					continue
				}

				ml.AddQuietMove(NewMove(sq, tSq, EMPTY, EMPTY, 0))
			}
		}

		piece = LoopNonSlidePieces[pieceIndex]
		pieceIndex++
	}
}
