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

	// iterative deepening
	for depth := 1; depth <= ss.MaxDepth; depth++ {
		if ss.ShouldStop(Soft) || ss.ShouldStop(Hard) {
			break
		}

		score := ss.negamax(pos, depth, 0, -INF, INF)

		if len(ss.pv.Line()) == 0 {
			// interruped before first move search completed,
			// use bestMove from previous depth
			break
		}

		bestMove = ss.pv.Line()[0]
		ss.printPV(os.Stdout, depth, score)
	}
	return bestMove
}
