package engine

/*
	Principal Variation Table
*/

// The maximum size of PvEntry on 64-bit systems is 16 bytes
const PVTABLESIZE = 0x100000 * 2 / 16

type PvEntry struct {
	PosKey uint64
	Move   int
}

type PvTable struct {
	Table [PVTABLESIZE]PvEntry
}

func (pv *PvTable) ClearPvTable() {
	for i := range pv.Table {
		pv.Table[i].Move = NOMOVE
		pv.Table[i].PosKey = 0
	}
}

func (pv *PvTable) StorePvMove(b *Board, move int) {
	index := b.PosKey % PVTABLESIZE

	pv.Table[index] = PvEntry{
		Move:   move,
		PosKey: b.PosKey,
	}
}

func (pv *PvTable) ProbePvTable(b *Board) int {
	index := b.PosKey % PVTABLESIZE
	entry := pv.Table[index]

	return entry.Move
}

func GetPvLine(depth int, b *Board) int {
	move := b.PvTable.ProbePvTable(b)
	count := 0
	for move != NOMOVE && count < depth {
		if MoveExists(b, move) == TRUE {
			b.MakeMove(move)
			b.PvArray[count] = move
			count++
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
