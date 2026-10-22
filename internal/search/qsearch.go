package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/tt"
)

func (ss *SearchState) qsearch(pos board.Position, ply uint16, alpha, beta int32) int32 {
	ss.pv.Clear(ply)
	ss.selDepth = max(ss.selDepth, ply)

	if ss.isDraw(&pos) {
		return 0
	}

	isPV := beta > alpha+1
	inCheck := pos.InCheck()

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	ttMove := board.NullMove
	if ttHit {
		score := scoreFromTT(entry.Score, ply)
		if entry.Flag.BoundsMatch(score, alpha, beta) {
			return score
		}
		if inCheck || entry.Move.IsNoisy() {
			ttMove = entry.Move
		}
	}

	bestScore := -INF
	bestMove := board.NullMove
	flag := tt.Upper

	if !inCheck {
		bestScore = evaluate(&pos) + ss.corrHist.Correction(&pos)
		if bestScore >= beta {
			return bestScore
		}
		alpha = max(alpha, bestScore)
	}

	movePicker := NewMovePicker(ttMove, !inCheck)
	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }

	legalMoves := 0
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		legalMoves++

		// SEE pruning
		if !inCheck && !see(&pos, move, 0) {
			continue
		}

		newPos := pos.MakeMove(move)
		ss.tt.Prefetch(newPos.Hash)
		ss.Nodes++

		ss.Keys.Push(newPos.Hash)
		score := -ss.qsearch(newPos, ply+1, -beta, -alpha)
		ss.Keys.Pop()

		if ss.ShouldStop(Hard) {
			return 0
		}

		bestScore = max(bestScore, score)

		if score > alpha {
			alpha = score
			bestMove = move
			flag = tt.Exact

			if isPV {
				ss.pv.Store(ply, move)
			}
			if score >= beta {
				flag = tt.Lower
				break
			}
		}
	}

	if inCheck && legalMoves == 0 {
		return matedIn(ply)
	}

	ss.tt.Store(pos.Hash, bestMove, scoreToTT(bestScore, ply), 0, flag)

	return bestScore
}
