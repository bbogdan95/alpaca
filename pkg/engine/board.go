package engine

import (
	"errors"
	"fmt"
	"io"
	"log"
)

type Undo struct {
	Move       int
	CastlePerm int
	EnPassant  int
	FiftyMove  int
	PosKey     uint64
}

type Board struct {
	Pieces    [BRD_SQ_NUM]int
	Pawns     [3]uint64
	KingSq    [2]int
	Side      int
	EnPassant int
	FiftyMove int
	Ply       int
	HisPly    int
	PosKey    uint64
	PCENum    [13]int
	BigPCE    [2]int
	MajPCE    [2]int
	MinPCE    [2]int
	Material  [2]int

	CastlePerm int
	History    [MAXGAMESMOVES]Undo
	PList      [13][10]int

	PvTable PvTable
	PvArray [MAXDEPTH]int
}

func (b *Board) PrintBoard(out io.Writer) {

	fmt.Fprintln(out, "Game board: ")
	fmt.Fprintln(out, "")

	for rank := RANK_8; rank >= RANK_1; rank-- {
		fmt.Fprintf(out, "%d", rank+1)
		for file := FILE_A; file <= FILE_H; file++ {
			sq := FR2SQ(file, rank)
			piece := b.Pieces[sq]
			fmt.Fprintf(out, " %3c", PceChar[piece])
		}
		fmt.Fprintln(out, "")
	}

	fmt.Fprintf(out, "\n ")
	for file := FILE_A; file <= FILE_H; file++ {
		fmt.Fprintf(out, " %3c", 'a'+file)
	}
	fmt.Fprintln(out, "")

	fmt.Fprintf(out, "side: %c\n", SideChar[b.Side])
	fmt.Fprintf(out, "EnPassant:%d\n", b.EnPassant)

	fmt.Fprintf(out, "castle: ")
	if b.CastlePerm&WKCA != 0 {
		fmt.Fprintf(out, "K")
	} else {
		fmt.Fprintf(out, "-")
	}

	if b.CastlePerm&WQCA != 0 {
		fmt.Fprintf(out, "Q")
	} else {
		fmt.Fprintf(out, "-")
	}

	if b.CastlePerm&BKCA != 0 {
		fmt.Fprintf(out, "k")
	} else {
		fmt.Fprintf(out, "-")
	}

	if b.CastlePerm&BQCA != 0 {
		fmt.Fprintf(out, "q")
	} else {
		fmt.Fprintf(out, "-")
	}
	fmt.Fprintf(out, "\n")

	fmt.Fprintf(out, "posKey: %X\n", b.PosKey)
}

func InitSq120To64() {
	sq := A1
	sq64 := 0

	for i := 0; i < BRD_SQ_NUM; i++ {
		SQ64[i] = 65
	}

	for i := 0; i < 64; i++ {
		SQ120[i] = 120
	}

	for rank := RANK_1; rank <= RANK_8; rank++ {
		for file := FILE_A; file <= FILE_H; file++ {
			sq = FR2SQ(file, rank)
			SQ120[sq64] = sq
			SQ64[sq] = sq64
			sq64++
		}
	}
}

func (b *Board) CheckBoard() {
	tempPceNum := [13]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	tempBigPce := [2]int{0, 0}
	tempMajPce := [2]int{0, 0}
	tempMinPce := [2]int{0, 0}
	tempMaterial := [2]int{0, 0}

	tempPawns := [3]uint64{0, 0, 0}

	tempPawns[WHITE] = b.Pawns[WHITE]
	tempPawns[BLACK] = b.Pawns[BLACK]
	tempPawns[BOTH] = b.Pawns[BOTH]

	for tempPiece := WP; tempPiece <= BK; tempPiece++ {
		for i := 0; i < b.PCENum[tempPiece]; i++ {
			sq120 := b.PList[tempPiece][i]
			if b.Pieces[sq120] != tempPiece {
				panic("Pieces not aligned (1)")
			}

		}
	}

	for sq64 := 0; sq64 < 64; sq64++ {
		sq120 := SQ120[sq64]
		tempPiece := b.Pieces[sq120]
		tempPceNum[tempPiece]++
		color := PieceCol[tempPiece]

		if PieceBig[tempPiece] == TRUE {
			tempBigPce[color]++
		}
		if PieceMaj[tempPiece] == TRUE {
			tempMajPce[color]++
		}
		if PieceMin[tempPiece] == TRUE {
			tempMinPce[color]++
		}

		if tempPiece != EMPTY {
			tempMaterial[color] += PieceVal[tempPiece]
		}
	}

	for tempPiece := WP; tempPiece <= BK; tempPiece++ {
		if tempPceNum[tempPiece] != b.PCENum[tempPiece] {
			panic("Pieces not aligned (2)")
		}
	}

	pcount := CountBits(tempPawns[WHITE])
	if pcount != b.PCENum[WP] {
		log.Fatalf("Pieces not aligned (3)")
	}

	pcount = CountBits(tempPawns[BLACK])
	if pcount != b.PCENum[BP] {
		log.Fatalf("Pieces not aligned (4)")
	}

	pcount = CountBits(tempPawns[BOTH])
	if pcount != b.PCENum[WP]+b.PCENum[BP] {
		log.Fatalf("Pieces not aligned (5)")
	}

}

func (b *Board) ResetBoard() {
	for i := 0; i < BRD_SQ_NUM; i++ {
		b.Pieces[i] = OFFBOARD
	}

	for i := 0; i < 64; i++ {
		b.Pieces[SQ120[i]] = EMPTY
	}

	for i := 0; i < 2; i++ {
		b.BigPCE[i] = 0
		b.MajPCE[i] = 0
		b.MinPCE[i] = 0
		b.Pawns[i] = 0
		b.Material[i] = 0
	}

	for i := 0; i < 3; i++ {
		b.Pawns[i] = 0
	}

	for i := 0; i < 13; i++ {
		b.PCENum[i] = 0
	}

	b.KingSq[WHITE] = NO_SQ
	b.KingSq[BLACK] = NO_SQ

	b.Side = BOTH
	b.EnPassant = NO_SQ
	b.FiftyMove = 0

	b.Ply = 0
	b.HisPly = 0
	b.CastlePerm = 0
	b.PosKey = 0

	b.PvTable = make(PvTable, 0)
}

func (b *Board) HashPiece(piece, sq int) {
	b.PosKey ^= PieceKeys[piece][sq]
}

func (b *Board) HashCA() {
	b.PosKey ^= CastleKeys[b.CastlePerm]
}

func (b *Board) HashSide() {
	b.PosKey ^= SideKey
}

func (b *Board) HashEP() {
	b.PosKey ^= PieceKeys[EMPTY][b.EnPassant]
}

func (b *Board) ClearPiece(sq int) {
	if SqOffBoard(sq) {
		panic("sq offboard")
	}

	b.CheckBoard()

	piece := b.Pieces[sq]
	col := PieceCol[piece]
	tPieceNum := -1

	if !PieceValid(piece) {
		panic("piece invalid")
	}

	b.HashPiece(piece, sq)
	b.Pieces[sq] = EMPTY
	b.Material[col] -= PieceVal[piece]

	if PieceBig[piece] != 0 {
		b.BigPCE[col]--
		if PieceMaj[piece] != 0 {
			b.MajPCE[col]--
		} else {
			b.MinPCE[col]--
		}
	} else {
		ClearBit(&b.Pawns[col], SQ64[sq])
		ClearBit(&b.Pawns[BOTH], SQ64[sq])
	}

	/*
		Assuming we have 5 white pawns on the board
		b.PCENum[WP] == 5 Looping from 0 to 4
		b.PList[piece][0] == sq0
		b.PList[piece][1] == sq1
		b.PList[piece][2] == sq2
		b.PList[piece][3] == sq3
		b.PList[piece][4] == sq4

		we loop until we find the sq that our piece is on
		sq == sq3 so tPieceNum = 3

		after removing the piece from PCENum and PList, it will look like this:
		b.PCENum[WP] == 4 Looping from 0 to 3
		b.PList[piece][0] == sq0
		b.PList[piece][1] == sq1
		b.PList[piece][2] == sq2
		b.PList[piece][3] == sq4
	*/
	for i := 0; i < b.PCENum[piece]; i++ {
		if b.PList[piece][i] == sq {
			tPieceNum = i
			break
		}
	}

	if tPieceNum == -1 {
		panic("couldn't find piece")
	}

	b.PCENum[piece]--
	b.PList[piece][tPieceNum] = b.PList[piece][b.PCENum[piece]]
}

func (b *Board) AddPiece(sq int, piece int) {
	if !PieceValid(piece) || SqOffBoard(sq) {
		panic("cannot add piece")
	}

	col := PieceCol[piece]
	b.HashPiece(piece, sq)
	b.Pieces[sq] = piece

	if PieceBig[piece] != 0 {
		b.BigPCE[col]++
		if PieceMaj[piece] != 0 {
			b.MajPCE[col]++
		} else {
			b.MinPCE[col]++
		}
	} else {
		SetBit(&b.Pawns[col], SQ64[sq])
		SetBit(&b.Pawns[BOTH], SQ64[sq])
	}

	b.Material[col] += PieceVal[piece]
	b.PList[piece][b.PCENum[piece]] = sq
	b.PCENum[piece]++
}

func (b *Board) MovePiece(from, to int) {
	if SqOffBoard(from) || SqOffBoard(to) {
		panic("cannot move piece offboard")
	}

	piece := b.Pieces[from]
	col := PieceCol[piece]

	b.HashPiece(piece, from)
	b.Pieces[from] = EMPTY

	b.HashPiece(piece, to)
	b.Pieces[to] = piece

	if PieceBig[piece] == 0 {
		ClearBit(&b.Pawns[col], SQ64[from])
		ClearBit(&b.Pawns[BOTH], SQ64[from])
		SetBit(&b.Pawns[col], SQ64[to])
		SetBit(&b.Pawns[BOTH], SQ64[to])
	}

	for i := 0; i < b.PCENum[piece]; i++ {
		if b.PList[piece][i] == from {
			b.PList[piece][i] = to
			break
		}
	}

}

func (b *Board) MakeMove(move int) (int, error) {
	b.CheckBoard()

	from := GetFrom(move)
	to := GetToSq(move)
	side := b.Side

	if SqOffBoard(from) || SqOffBoard(to) || !SideValid(side) || !PieceValid(b.Pieces[from]) {

		return 0, errors.New("cannot make move")
	}

	b.History[b.HisPly].PosKey = b.PosKey

	if move&MoveFlagEnPassant != 0 {
		if side == WHITE {
			b.ClearPiece(to - 10)
		} else {
			b.ClearPiece(to + 10)
		}
	} else if move&MoveFlagCastle != 0 {
		switch to {
		case C1:
			b.MovePiece(A1, D1)
		case C8:
			b.MovePiece(A8, D8)
		case G1:
			b.MovePiece(H1, F1)
		case G8:
			b.MovePiece(H8, F8)
		default:
			return 0, errors.New("illegal castle move")
		}
	}

	if b.EnPassant != NO_SQ {
		b.HashEP()
	}
	b.HashCA()

	b.History[b.HisPly].Move = move
	b.History[b.HisPly].FiftyMove = b.FiftyMove
	b.History[b.HisPly].EnPassant = b.EnPassant
	b.History[b.HisPly].CastlePerm = b.CastlePerm

	b.CastlePerm &= CastlePerm[from]
	b.CastlePerm &= CastlePerm[to]
	b.EnPassant = NO_SQ

	b.HashCA()

	captured := GetCaptured(move)
	b.FiftyMove++

	if captured != EMPTY {
		if !PieceValid(captured) {

			return 0, errors.New("captured piece invalid")
		}

		b.ClearPiece(to)
		b.FiftyMove = 0
	}

	b.HisPly++
	b.Ply++

	if PiecePawn[b.Pieces[from]] != 0 {
		b.FiftyMove = 0
		if move&MoveFlagPawnStart != 0 {
			if side == WHITE {
				b.EnPassant = from + 10
				if RanksBrd[b.EnPassant] != RANK_3 {
					return 0, errors.New("illegal pawn start move")
				}
			} else {
				b.EnPassant = from - 10
				if RanksBrd[b.EnPassant] != RANK_6 {
					return 0, errors.New("illegal pawn start move")
				}
			}
			b.HashEP()
		}
	}

	b.MovePiece(from, to)

	isPromotedPiece := GetPromoted(move)
	if isPromotedPiece != EMPTY {
		if !PieceValid(isPromotedPiece) || PiecePawn[isPromotedPiece] == TRUE {
			return 0, errors.New("illegal piece promotion")
		}

		b.ClearPiece(to)
		b.AddPiece(to, isPromotedPiece)
	}

	if PieceKing[b.Pieces[to]] != 0 {
		b.KingSq[b.Side] = to
	}

	b.Side ^= 1
	b.HashSide()

	b.CheckBoard()

	if SqAttacked(b.KingSq[side], b.Side, b) == TRUE {
		b.TakeMove()
		return FALSE, nil
	}

	return TRUE, nil
}

func (b *Board) TakeMove() {

	b.CheckBoard()

	b.HisPly--
	b.Ply--

	move := b.History[b.HisPly].Move
	from := GetFrom(move)
	to := GetToSq(move)

	if SqOffBoard(from) || SqOffBoard(to) {
		panic("can't take move - sq offboard")
	}

	if b.EnPassant != NO_SQ {
		b.HashEP()
	}
	b.HashCA()

	b.CastlePerm = b.History[b.HisPly].CastlePerm
	b.FiftyMove = b.History[b.HisPly].FiftyMove
	b.EnPassant = b.History[b.HisPly].EnPassant

	if b.EnPassant != NO_SQ {
		b.HashEP()
	}
	b.HashCA()

	b.Side ^= 1
	b.HashSide()

	if MoveFlagEnPassant&move != 0 {
		if b.Side == WHITE {
			b.AddPiece(to-10, BP)
		} else {
			b.AddPiece(to+10, WP)
		}
	} else if MoveFlagCastle&move != 0 {
		switch to {
		case C1:
			b.MovePiece(D1, A1)
		case C8:
			b.MovePiece(D8, A8)
		case G1:
			b.MovePiece(F1, H1)
		case G8:
			b.MovePiece(F8, H8)
		default:
			panic("illegal take move")
		}
	}

	b.MovePiece(to, from)

	if PieceKing[b.Pieces[from]] == TRUE {
		b.KingSq[b.Side] = from
	}

	captured := GetCaptured(move)
	if captured != EMPTY {
		if !PieceValid(captured) {
			panic("invalid captured piece in take move")
		}

		b.AddPiece(to, captured)
	}

	promoted := GetPromoted(move)
	if promoted != EMPTY {
		if !PieceValid(promoted) || PiecePawn[promoted] == TRUE {
			panic("illegal promoted piece in take move")
		}

		b.ClearPiece(from)
		piece := WP
		if PieceCol[promoted] == BLACK {
			piece = BP
		}
		b.AddPiece(from, piece)
	}

	b.CheckBoard()
}

// from the last time the FiftyMove was set to zero, loop over and check for repetition
func (b *Board) IsRepetition() bool {
	for i := b.HisPly - b.FiftyMove; i < b.HisPly-1; i++ {
		// will we ever go over MAXGAMEMOVES ?
		if b.PosKey == b.History[i].PosKey {
			return true
		}
	}

	return false
}
