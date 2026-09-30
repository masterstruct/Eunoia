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
	bestMove := firstLegalMove(&pos)
	if bestMove == board.NullMove {
		return board.NullMove
	}

	var lastScore int16
	aw := newAspirationWindow() // [-INF; +INF]

	// iterative deepening
iterativeDeepening:
	for depth := 1; depth <= ss.MaxDepth; depth++ {
		if ss.ShouldStop(Soft) || ss.ShouldStop(Hard) {
			break
		}

		if depth > 5 {
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
				aw.alpha = -INF
				aw.beta = INF
				// aw.widenDown()
				continue
			}
			if score >= aw.beta {
				aw.alpha = -INF
				aw.beta = INF
				// aw.widenUp()
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

func firstLegalMove(pos *board.Position) board.Move {
	var movelist board.Movelist
	board.GenerateLegalMoves(pos, &movelist, board.All)
	if movelist.Len > 0 {
		return movelist.Moves[0]
	}
	return board.NullMove
}
