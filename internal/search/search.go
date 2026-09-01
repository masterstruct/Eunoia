package search

import (
	"os"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
)

const (
	MATE int16 = 30000
	INF  int16 = 32000
)

func (ss *SearchState) SearchBestMove(pos board.Position, maxDepth int) board.Move {
	var bestMove board.Move

	var lastScore int16
	aw := newAspirationWindow() // [-INF; +INF]

	// iterative deepening
iterativeDeepening:
	for depth := 1; depth <= maxDepth; depth++ {
		if (ss.SoftNodes > 0 && ss.Nodes >= ss.SoftNodes) || (!ss.SoftTime.IsZero() && time.Now().After(ss.SoftTime)) {
			break
		}

		if depth > 1 {
			aw.centerAround(lastScore)
		}

		// aspiration search
		for {
			score := ss.negamax(pos, depth, 0, aw.alpha, aw.beta)

			if ss.searchStopped() {
				break iterativeDeepening
			}

			// TODO: widen gradually
			if score <= aw.alpha {
				// aw.widenDown()
				aw.alpha = -INF
				continue
			}
			if score >= aw.beta {
				// aw.widenUp()
				aw.beta = INF
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

func (ss *SearchState) searchStopped() bool {
	if ss.Stop {
		return true
	}
	if ss.MaxNodes > 0 && ss.Nodes >= ss.MaxNodes {
		return true
	}
	if !ss.MaxTime.IsZero() && time.Now().After(ss.MaxTime) {
		return true
	}
	return false
}
