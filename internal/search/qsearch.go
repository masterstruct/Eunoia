package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func (ss *SearchState) qsearch(pos board.Position, ply uint16, alpha, beta int32) int32 {
	ss.pv.Clear(ply)
	ss.selDepth = max(ss.selDepth, ply)

	if ss.isDraw(&pos) {
		return 0
	}

	isPV := beta > alpha+1
	inCheck := pos.InCheck()

	best := -INF
	if !inCheck {
		standPat := evaluate(&pos) + ss.corrHist.Correction(&pos)
		if standPat >= beta {
			return standPat
		}
		best = standPat
		alpha = max(alpha, standPat)
	}

	movePicker := NewMovePicker(ply, board.NullMove, !inCheck)
	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }

	legalMoves := 0
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		legalMoves++

		// SEE pruning
		if !inCheck && !see(&pos, move, 0) {
			continue
		}

		newPos := pos
		ss.makeMove(&newPos, move)
		score := -ss.qsearch(newPos, ply+1, -beta, -alpha)
		ss.unmakeMove()

		if ss.ShouldStop(Hard) {
			return 0
		}

		if score > best {
			best = score
			if score > alpha {
				alpha = score
				if isPV {
					ss.pv.Store(ply, move)
				}
				if score >= beta {
					break
				}
			}
		}
	}

	if inCheck && legalMoves == 0 {
		return matedIn(ply)
	}
	return best
}
