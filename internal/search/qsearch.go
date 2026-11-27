package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/history"
)

const (
	qsFutilityMargin = 150
)

func (ss *SearchState) qsearch(pos board.Position, ply uint16, alpha, beta int32) int32 {
	ss.pv.Clear(ply)
	ss.selDepth = max(ss.selDepth, ply)

	if ss.isDraw(&pos) {
		return 0
	}

	if ply >= MaxPly {
		return evaluate(&pos)
	}

	isPV := beta > alpha+1
	inCheck := pos.InCheck()

	staticEval := -INF
	if !inCheck {
		staticEval = clampMate(evaluate(&pos) + ss.corrHist.Correction(&pos))
		if staticEval >= beta {
			return staticEval
		}
		alpha = max(alpha, staticEval)
	}

	movePicker := NewMovePicker(ply, board.NullMove, !inCheck, history.ContIndices{})
	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }

	futility := staticEval + qsFutilityMargin
	bestScore := staticEval

	legalMoves := 0
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		legalMoves++
		if !isMated(bestScore) {
			// futility pruning
			if futility <= alpha && !see(&pos, move, 1) {
				bestScore = max(bestScore, futility)
				continue
			}

			// SEE pruning
			if !inCheck && !see(&pos, move, 0) {
				continue
			}
		}

		newPos := pos
		ss.makeMove(&newPos, move, ply)
		score := -ss.qsearch(newPos, ply+1, -beta, -alpha)
		ss.unmakeMove(ply)

		if ss.ShouldStop(Hard) {
			return 0
		}

		if score > bestScore {
			bestScore = score
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
	return bestScore
}
