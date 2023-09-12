package engine

const (
	EMPTY = iota
	WP
	WN
	WB
	WR
	WQ
	WK
	BP
	BN
	BB
	BR
	BQ
	BK
)

const (
	FILE_A = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
	FILE_NONE
)

const (
	RANK_1 = iota
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
	RANK_NONE
)

const (
	WHITE = iota
	BLACK
	BOTH
)

const (
	A1 = 21 + iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
)

const (
	A2 = 31 + iota
	B2
	C2
	D2
	E2
	F2
	G2
	H2
)

const (
	A3 = 41 + iota
	B3
	C3
	D3
	E3
	F3
	G3
	H3
)

const (
	A4 = 51 + iota
	B4
	C4
	D4
	E4
	F4
	G4
	H4
)

const (
	A5 = 61 + iota
	B5
	C5
	D5
	E5
	F5
	G5
	H5
)

const (
	A6 = 71 + iota
	B6
	C6
	D6
	E6
	F6
	G6
	H6
)

const (
	A7 = 81 + iota
	B7
	C7
	D7
	E7
	F7
	G7
	H7
)

const (
	A8 = 91 + iota
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	NO_SQ
	OFFBOARD
)

const (
	BRD_SQ_NUM = 120
	WKCA       = 1
	WQCA       = 2
	BKCA       = 4
	BQCA       = 8

	MAXGAMESMOVES    = 2048
	MAXPOSITIONMOVES = 256
)

const (
	FALSE = iota
	TRUE
)
