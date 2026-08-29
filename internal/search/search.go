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

	// iterative deepening
	for depth := 1; depth <= maxDepth; depth++ {
		if (ss.SoftNodes > 0 && ss.Nodes >= ss.SoftNodes) || (!ss.SoftTime.IsZero() && time.Now().After(ss.SoftTime)) {
			break
		}

		score := ss.negamax(pos, depth, 0, -INF, INF)

		if len(ss.pv.Line()) == 0 {
			// interruped before first move search completed,
			// use bestMove from previous depth
			break
		}

		bestMove = ss.pv.Line()[0]
		ss.printPV(os.Stdout, depth, score, pos.SideToMove)

		if ss.searchStopped() {
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
