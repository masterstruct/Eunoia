package movegen

import (
	"fmt"

	"github.com/masterstruct/Eunoia/internal/board"
)

func Perft(pos board.Position, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var nodes uint64
	movelist := GeneratePseudolegalMoves(pos)
	for _, move := range movelist {
		mover := pos.SideToMove
		newPos := pos.MakeMove(move)
		if InCheck(newPos, mover) {
			continue
		}
		nodes += Perft(newPos, depth-1)
	}
	return nodes
}

func SplitPerft(pos board.Position, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var nodes uint64
	var total uint64
	movelist := GeneratePseudolegalMoves(pos)
	for _, move := range movelist {
		mover := pos.SideToMove
		newPos := pos.MakeMove(move)
		if InCheck(newPos, mover) {
			continue
		}
		nodes = Perft(newPos, depth-1)
		fmt.Println(move, nodes)
		total += nodes
	}
	fmt.Println("Total:", total)
	return total
}
