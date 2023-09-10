package perft

// import (
// 	"fmt"
// 	"os"
// 	"sync"
// )

// type tree struct {
// 	LeafNodes uint64
// 	Mu        sync.Mutex
// }

// // var leafNodes uint64
// var leafNodes tree

// func Perft(depth int, b *Board, wg *sync.WaitGroup) {
// 	b.CheckBoard()
// 	wg.Add(1)
// 	defer wg.Done()

// 	if depth == 0 {
// 		leafNodes.Mu.Lock()
// 		leafNodes.LeafNodes++
// 		leafNodes.Mu.Unlock()

// 		return
// 	}

// 	ml := &MoveList{}
// 	GenerateAllMoves(b, ml)

// 	for i := 0; i < ml.Count; i++ {
// 		res, err := b.MakeMove(ml.Moves[i].Move)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if res == FALSE {
// 			continue
// 		}

// 		Perft(depth-1, b, wg)
// 		b.TakeMove()
// 	}
// }

// func PerftTest(depth int, b *Board) {
// 	b.CheckBoard()

// 	b.PrintBoard(os.Stdout)
// 	fmt.Fprintf(os.Stdout, "\nStarting Test To Depth:%d\n", depth)

// 	leafNodes.LeafNodes = 0
// 	ml := &MoveList{}
// 	GenerateAllMoves(b, ml)

// 	var wg sync.WaitGroup

// 	for i := 0; i < ml.Count; i++ {
// 		res, err := b.MakeMove(ml.Moves[i].Move)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if res == FALSE {
// 			continue
// 		}

// 		leafNodes.Mu.Lock()
// 		cumnodes := leafNodes.LeafNodes
// 		leafNodes.Mu.Unlock()
// 		cb := *b
// 		go Perft(depth-1, &cb, &wg)
// 		b.TakeMove()
// 		leafNodes.Mu.Lock()
// 		oldnodes := leafNodes.LeafNodes - cumnodes
// 		leafNodes.Mu.Unlock()

// 		fmt.Fprintf(os.Stdout, "move %d : %s : %1d\n", i+1, PrintMove(ml.Moves[i].Move), oldnodes)
// 	}

// 	wg.Wait()

// 	leafNodes.Mu.Lock()
// 	fmt.Fprintf(os.Stdout, "\nTest Complete : %d nodes visited\n", leafNodes.LeafNodes)
// 	leafNodes.Mu.Unlock()
// }
//
