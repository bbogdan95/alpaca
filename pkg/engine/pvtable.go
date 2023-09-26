package engine

import "fmt"

const (
	HFNONE = iota
	HFALPHA
	HFBETA
	HFEXACT
)

type HashEntry struct {
	PosKey uint64
	Move   int
	Score  int
	Depth  int
	Flags  int
}

type HashTable struct {
	Table      map[string]HashEntry
	NumEntries int
	NewWrite   int
	OverWrite  int
	Hit        int
	Cut        int
}

func (ht *HashTable) ClearHashTable() {
	ht.Table = map[string]HashEntry{}
	ht.NewWrite = 0
}

func (ht *HashTable) StoreHashEntry(b *Board, move int, score int, flags int, depth int) {
	indexStr := fmt.Sprintf("%d", b.PosKey)

	if b.HashTable.Table[indexStr].PosKey == 0 {
		b.HashTable.NewWrite++
	} else {
		b.HashTable.OverWrite++
	}

	if score > ISMATE {
		score += b.Ply
	} else if score < -ISMATE {
		score -= b.Ply
	}

	ht.Table[indexStr] = HashEntry{
		Move:   move,
		PosKey: b.PosKey,
		Flags:  flags,
		Score:  score,
		Depth:  depth,
	}
}

func (ht *HashTable) ProbeHashEntry(b *Board, move *int, score *int, alpha int, beta int, depth int) int {
	indexStr := fmt.Sprintf("%d", b.PosKey)

	if ht.Table[indexStr].PosKey == b.PosKey {
		*move = ht.Table[indexStr].Move
		if ht.Table[indexStr].Depth >= depth {
			ht.Hit++

			*score = ht.Table[indexStr].Score
			if *score > ISMATE {
				*score -= b.Ply
			} else {
				if *score < -ISMATE {
					*score += b.Ply
				}
			}

			switch ht.Table[indexStr].Flags {
			case HFALPHA:
				if *score <= alpha {
					*score = alpha
					return TRUE
				}
			case HFBETA:
				if *score >= beta {
					*score = beta
					return TRUE
				}
			case HFEXACT:
				return TRUE
			default:
				panic(false)
			}
		}
	}

	return FALSE
}

func GetPvLine(depth int, b *Board) int {
	move := b.HashTable.ProbePvMove(b)
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

		move = b.HashTable.ProbePvMove(b)
	}

	for b.Ply > 0 {
		b.TakeMove()
	}

	return count
}

func (ht *HashTable) ProbePvMove(b *Board) int {
	indexStr := fmt.Sprintf("%d", b.PosKey)

	if ht.Table[indexStr].PosKey == b.PosKey {
		return ht.Table[indexStr].Move
	}

	return NOMOVE
}
