package engine

import (
	"unsafe"
)

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
	Table      []HashEntry
	NumEntries int
	NewWrite   int
	OverWrite  int
	Hit        int
	Cut        int
}

// HashEntry has a size of 40 bytes on x64 systems.
// In order to set a limit to our HashTable in terms of MB, the internal representation of the HashTable
// is a []HashEntry with a preallocated capacity of 1024 * 1024 * MB (size in bytes). We calculate the maximum number of entries
// in our HashTable by taking the desired size in bytes (1024 * 1024 * MB) and dividing it by the size of one HashEntry (40 bytes)
// Our HashTable will never grow beyound the predefined limit (in terms of MB but) - but it will also never grow beyond b.HashTable.NumEntries
// because when we store something into it, we use the zobrist key (uint64) % number of elements - the result of this operation is always
// a number between 0 and b.HashTable.NumEntries.
func InitHashTable(b *Board, MB int) {
	entrySize := unsafe.Sizeof(HashEntry{}) // 40 bytes

	HashSize := 1024 * 1024 * MB
	b.HashTable.NumEntries = HashSize / int(entrySize)

	b.HashTable.Table = make([]HashEntry, b.HashTable.NumEntries)
}

func ClearHashTable(b *Board) {
	b.HashTable.Table = make([]HashEntry, b.HashTable.NumEntries)
	b.HashTable.NewWrite = 0
}

func StoreHashEntry(b *Board, move int, score int, flags int, depth int) {
	// by doing this, we limit the number of entries to our defined size in MB
	index := b.PosKey % uint64(b.HashTable.NumEntries)
	if b.HashTable.Table[index].PosKey == 0 {
		b.HashTable.NewWrite++
	} else {
		b.HashTable.OverWrite++
	}

	if score > ISMATE {
		score += b.Ply
	} else if score < -ISMATE {
		score -= b.Ply
	}

	b.HashTable.Table[index] = HashEntry{
		Move:   move,
		PosKey: b.PosKey,
		Flags:  flags,
		Score:  score,
		Depth:  depth,
	}
}

func ProbeHashEntry(b *Board, move *int, score *int, alpha int, beta int, depth int) int {
	index := b.PosKey % uint64(b.HashTable.NumEntries)
	if b.HashTable.Table[index].PosKey == b.PosKey {
		*move = b.HashTable.Table[index].Move
		if b.HashTable.Table[index].Depth >= depth {
			b.HashTable.Hit++

			*score = b.HashTable.Table[index].Score
			if *score > ISMATE {
				*score -= b.Ply
			} else {
				if *score < -ISMATE {
					*score += b.Ply
				}
			}

			switch b.HashTable.Table[index].Flags {
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
	move := ProbePvMove(b)
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

		move = ProbePvMove(b)
	}

	for b.Ply > 0 {
		b.TakeMove()
	}

	return count
}

func ProbePvMove(b *Board) int {
	index := b.PosKey % uint64(b.HashTable.NumEntries)
	if b.HashTable.Table[index].PosKey == b.PosKey {
		return b.HashTable.Table[index].Move
	}

	return NOMOVE
}
