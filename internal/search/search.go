package search

import (
	"os"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/tt"
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
	for depth := 1; depth <= int(ss.MaxDepth); depth++ {
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
				ss.pv.EnsureRoot(bestMove)
				ss.printPV(os.Stdout, depth, score, tt.Upper)

				aw.alpha = -INF
				aw.beta = INF
				// aw.widenDown()
				continue
			}
			bestMove = ss.pv.BestMove()
			if score >= aw.beta {
				ss.printPV(os.Stdout, depth, score, tt.Lower)

				aw.alpha = -INF
				aw.beta = INF
				// aw.widenUp()
				continue
			}

			ss.printPV(os.Stdout, depth, score, tt.Exact)
			lastScore = score
			break
		}
	}
	return bestMove
}

func firstLegalMove(pos *board.Position) board.Move {
	var movelist board.Movelist
	board.GenerateLegalMoves(pos, &movelist, board.All)
	return movelist.Moves[0].Move // NullMove if no legal moves
}
