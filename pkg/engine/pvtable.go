package engine

import "fmt"

/*
	Principal Variation Table
*/
type PvEntry struct {
	PosKey uint64
	Move   int
}

type PvTable struct {
	Table map[string]PvEntry
}

func (pv *PvTable) ClearPvTable() {
	pv.Table = map[string]PvEntry{}
}

func (pv *PvTable) StorePvMove(b *Board, move int) {
	indexStr := fmt.Sprintf("%d", b.PosKey)

	pv.Table[indexStr] = PvEntry{
		Move:   move,
		PosKey: b.PosKey,
	}
}

func (pv *PvTable) ProbePvTable(b *Board) int {
	indexStr := fmt.Sprintf("%d", b.PosKey)
	entry := pv.Table[indexStr]

	return entry.Move
}

func GetPvLine(depth int, b *Board) int {
	move := b.PvTable.ProbePvTable(b)
	count := 0

	for move != NOMOVE && count < depth {
		c, err := MoveExists(b, move)
		if err != nil {
			panic(err)
		}
		if c == TRUE {
			res, err := b.MakeMove(move)
			if err != nil {
				panic(err)
			}
			if res == TRUE {
				b.PvArray[count] = move
				count++
			}
		} else {
			break
		}

		move = b.PvTable.ProbePvTable(b)
	}

	for b.Ply > 0 {
		b.TakeMove()
	}

	return count
}
