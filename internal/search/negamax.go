package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/tt"
)

func (ss *SearchState) negamax(pos board.Position, depth, ply int, alpha, beta int16) int16 {
	ss.pv.Init(ply)

	if ss.searchStopped() {
		return 0
	}

	ss.Nodes++

	alphaOrig := alpha

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	if ttHit && ply > 0 && entry.Depth >= uint8(depth) {
		score := scoreFromTT(entry.Score, ply)
		switch entry.Flag {
		case tt.Exact:
			return score
		case tt.Lower:
			if score >= beta {
				return score
			}
		case tt.Upper:
			if score <= alpha {
				return score
			}
		}
	}

	if depth <= 0 {
		return ss.qsearch(&pos, alpha, beta)
	}

	bestValue := -INF
	var bestMove board.Move
	legalMoves := 0

	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(&pos, &movelist)
	ss.orderMoves(&pos, &movelist)
	mover := pos.SideToMove

	for i := range movelist.Len {
		move := movelist.Moves[i]
		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}

		legalMoves++

		score := -ss.negamax(newPos, depth-1, ply+1, -beta, -alpha)

		if ss.searchStopped() {
			return 0
		}

		if score > bestValue {
			bestValue = score
			bestMove = move
			if score > alpha {
				alpha = score
				ss.pv.Store(ply, move)
			}
		}
		if score >= beta {
			break
		}
	}

	if legalMoves == 0 {
		if movegen.InCheck(&pos, mover) {
			// checkmate
			return -MATE + int16(ply)
		}
		// stalemate
		return 0
	}

	flag := tt.Exact
	if bestValue <= alphaOrig {
		flag = tt.Upper
	} else if bestValue >= beta {
		flag = tt.Lower
	}

	ss.tt.Store(pos.Hash, bestMove, scoreToTT(bestValue, ply), uint8(depth), flag)

	return bestValue
}

func scoreToTT(score int16, ply int) int16 {
	if score >= MATE-int16(MaxPly) {
		return score + int16(ply)
	}
	if score <= -MATE+int16(MaxPly) {
		return score - int16(ply)
	}
	return score
}

func scoreFromTT(score int16, ply int) int16 {
	if score >= MATE-int16(MaxPly) {
		return score - int16(ply)
	}
	if score <= -MATE+int16(MaxPly) {
		return score + int16(ply)
	}
	return score
}
