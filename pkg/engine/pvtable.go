package engine

import "fmt"

/*
	Principal Variation Table
*/

type PvTable []int

func (table *PvTable) ClearPvTable() {
	*table = []int{}
}

func (table *PvTable) StorePvMove(move int) {
	*table = append([]int{move}, *table...)
}

func (table *PvTable) ProbePvTable() int {
	if len(*table) > 0 {
		return (*table)[0]
	}

	return NOMOVE
}

func GetPvLine(depth int, b *Board) int {
	move := b.PvTable.ProbePvTable()
	fmt.Println("getpv line - ", PrintMove(move))
	count := 0
	for move != NOMOVE && count < depth {
		if MoveExists(b, move) == TRUE {
			fmt.Println(111111111111)
			b.MakeMove(move)
			count++
			b.PvArray[count] = move
		} else {
			fmt.Println(222222222222)
			break
		}

		move = b.PvTable.ProbePvTable()
	}

	for b.Ply > 0 {
		b.TakeMove()
	}

	return count
}
