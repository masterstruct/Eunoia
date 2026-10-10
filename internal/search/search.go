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
	switch rootMoves.Len {
	case 0:
		return ss.handleNoLegalMoves(&pos)
	case 1:
		return ss.handleOneLegalMove(&pos, &rootMoves)
	}

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
			if ss.depth > 1 && ss.Stop {
				break iterativeDeepening
			}

			bound := tt.FlagFromScore(score, aw.alpha, aw.beta)
			ss.printPV(os.Stdout, max(aw.alpha, min(score, aw.beta)), bound)

			if isMateScore(score) {
				ss.bestScore = score
				break
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
			ss.bestScore = score
			break
		}

		ss.completedDepth = ss.depth
		ss.depth++
	}

	ss.printInfo(os.Stdout, ss.completedDepth, ss.bestScore, tt.Exact, ss.rootPV[:ss.rootPVLen])

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

func (ss *SearchState) handleOneLegalMove(pos *board.Position, rootMoves *board.Movelist) board.Move {
	move := rootMoves.Moves[0].Move
	ss.bestMove = move
	ss.bestScore = evaluate(pos)

	ss.pv.Clear(1)
	ss.pv.Store(0, move)

	ss.printPV(os.Stdout, ss.bestScore, tt.Exact)
	return move
}
