package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/tt"
)

func (ss *SearchState) qsearch(pos *board.Position, ply int, alpha, beta int16) int16 {
	ss.Nodes++

	isPV := beta > alpha+1

	// TT lookup
	// TODO: tt.Store?
	entry, ttHit := ss.tt.Probe(pos.Hash)
	if ttHit && !isPV {
		score := scoreFromTT(entry.Score, ply)
		switch entry.Flag {
		case tt.Exact:
			return score
		case tt.Lower:
			if score >= beta {
				return score
			}
		case tt.Upper:
			if score <= alpha {
				return score
			}
		}
	}

	standPat := evaluate(pos)

	if standPat >= beta {
		return standPat
	}
	if standPat > alpha {
		alpha = standPat
	}

	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(pos, &movelist)
	ss.orderMoves(pos, &movelist, entry.Move)
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

		score := -ss.qsearch(&newPos, ply+1, -beta, -alpha)

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
