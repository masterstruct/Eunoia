package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
)

func (ss *SearchState) qsearch(pos *board.Position, alpha, beta int16) int16 {
	ss.Nodes++

	standPat := evaluate(pos)

	if standPat >= beta {
		return standPat
	}
	if standPat > alpha {
		alpha = standPat
	}

	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(pos, &movelist)
	ss.orderMoves(pos, &movelist)
	mover := pos.SideToMove

	for i := range movelist.Len {
		move := movelist.Moves[i]
		if !move.IsCapture() && !move.IsPromo() {
			continue
		}

		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}

		score := -ss.qsearch(&newPos, -beta, -alpha)

		if ss.searchStopped() {
			return 0
		}

		if score >= beta {
			return score
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
