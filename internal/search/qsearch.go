package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func (ss *SearchState) qsearch(pos board.Position, alpha, beta int16) int16 {
	ss.Nodes++

	standPat := evaluate(&pos)

	if standPat >= beta {
		return standPat
	}
	if standPat > alpha {
		alpha = standPat
	}

	var movelist board.Movelist
	board.GenerateLegalMoves(&pos, &movelist, board.All)
	ss.orderMoves(&pos, &movelist)

	for i := range movelist.Len {
		move := movelist.Moves[i]
		if !move.IsCapture() && !move.IsPromo() {
			continue
		}

		newPos := pos.MakeMove(move)

		score := -ss.qsearch(newPos, -beta, -alpha)

		if ss.ShouldStop(Hard) {
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
