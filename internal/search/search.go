package search

import (
	"fmt"
	"os"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/tt"
)

func (ss *SearchState) SearchBestMove(pos board.Position) board.Move {
	ss.pv.Clear(0)

	var rootMoves board.Movelist
	board.GenerateLegalMoves(&pos, &rootMoves, board.All)
	if rootMoves.Len == 0 {
		return ss.handleNoLegalMoves(&pos)
	}

	ss.bestMove = rootMoves.Moves[0].Move

	aw := newAspirationWindow() // [-INF; +INF]
	var score int32

	// iterative deepening
iterativeDeepening:
	for ss.depth <= int(MaxPly) && !ss.ShouldStop(Soft) {
		if ss.depth >= aspirationMinDepth {
			aw.centerAround(score)
		}

		// aspiration search
		for {
			score = ss.negamax(pos, ss.depth, 0, aw.alpha, aw.beta)
			if ss.ShouldStop(Hard) {
				break iterativeDeepening
			}

			bound := tt.FlagFromScore(score, aw.alpha, aw.beta)
			if line := ss.pv.Line(); len(line) > 0 {
				ss.bestMove = line[0]
			}
			ss.pv.EnsureRoot(ss.bestMove)
			ss.printPV(os.Stdout, max(aw.alpha, min(score, aw.beta)), bound)

			if score <= aw.alpha {
				aw.widenDown(score)
				continue
			}
			if score >= aw.beta {
				aw.widenUp(score)
				continue
			}
			ss.bestScore = score
			break
		}

		ss.depth++
	}

	return ss.bestMove
}

func (ss *SearchState) handleNoLegalMoves(pos *board.Position) board.Move {
	fmt.Println("info error no legal moves")

	score := int32(0) // stalemate
	if pos.InCheck() {
		score = -MATE
	}

	ss.bestMove = board.NullMove
	ss.bestScore = score
	return ss.bestMove
}
