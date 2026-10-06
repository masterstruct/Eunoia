package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func (ss *SearchState) qsearch(pos board.Position, ply uint16, alpha, beta int16) int16 {
	ss.Nodes++
	inCheck := pos.InCheck()

	best := -INF
	if !inCheck {
		standPat := evaluate(&pos)
		if standPat >= beta {
			return standPat
		}
		best = standPat
		alpha = max(alpha, standPat)
	}

	movePicker := NewMovePicker(board.NullMove, !inCheck)
	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }

	legalMoves := 0
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		legalMoves++
		score := -ss.qsearch(pos.MakeMove(move), ply+1, -beta, -alpha)
		if ss.ShouldStop(Hard) {
			return 0
		}
		if score > best {
			best = score
			if score > alpha {
				alpha = score
				if score >= beta {
					break
				}
			}
		}
	}

	if inCheck && legalMoves == 0 {
		return -MATE + int16(ply)
	}
	return best
}
