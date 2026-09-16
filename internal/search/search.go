package search

import (
	"os"

	"github.com/masterstruct/Eunoia/internal/board"
)

const (
	MATE int16 = 30000
	INF  int16 = 32000
)

func (ss *SearchState) SearchBestMove(pos board.Position) board.Move {
	var bestMove board.Move

	var lastScore int16
	aw := newAspirationWindow() // [-INF; +INF]

	// iterative deepening
iterativeDeepening:
	for depth := 1; depth <= ss.MaxDepth; depth++ {
		if ss.ShouldStop(Soft) || ss.ShouldStop(Hard) {
			break
		}

		if depth >= aspirationMinDepth {
			aw.centerAround(lastScore)
		}

		// aspiration search
		for {
			score := ss.negamax(pos, depth, 0, aw.alpha, aw.beta)

			if ss.ShouldStop(Hard) {
				// interruped before first move search completed,
				// discard results from this depth
				break iterativeDeepening
			}

			if score <= aw.alpha {
				aw.widenDown(score)
				continue
			}
			if score >= aw.beta {
				aw.widenUp(score)
				continue
			}

			bestMove = ss.pv.Line()[0]
			ss.printPV(os.Stdout, depth, score)
			lastScore = score
			break
		}
	}
	return bestMove
}
