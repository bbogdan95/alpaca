package main

import "math/rand"

var PieceKeys [13][120]uint64
var SideKey uint64
var CastleKeys [16]uint64

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

func GeneratePosKey(b *Board) uint64 {
	var finalKey uint64 = 0
	piece := EMPTY

	for sq := 0; sq < BRD_SQ_NUM; sq++ {
		piece = b.Pieces[sq]
		if piece != NO_SQ && piece != EMPTY && piece != OFFBOARD {
			//if piece >= WP && piece <= BK {
			finalKey ^= PieceKeys[piece][sq]
			//}
		}
	}

	if b.Side == WHITE {
		finalKey ^= SideKey
	}

	if b.EnPassant != NO_SQ {
		//if b.EnPassant >= 0 && b.EnPassant < BRD_SQ_NUM {
		finalKey ^= PieceKeys[EMPTY][b.EnPassant]
		//}
	}

	//if b.CastlePerm >= 0 && b.CastlePerm <= 15 {
	finalKey ^= CastleKeys[b.CastlePerm]
	//}

	return finalKey
}
