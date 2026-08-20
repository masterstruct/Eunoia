package search

import (
	"math/rand/v2"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
)

func Search(pos board.Position) board.Move {
	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(&pos, &movelist)
	mover := pos.SideToMove

	order := make([]int, movelist.Len)
	for i := range order {
		order[i] = i
	}
	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	for _, idx := range order {
		move := movelist.Moves[idx]
		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}
		return move
	}

	return board.NullMove
}
