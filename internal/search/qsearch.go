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

	movePicker := NewMovePicker(board.NullMove, true)
	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		score := -ss.qsearch(pos.MakeMove(move), -beta, -alpha)

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
