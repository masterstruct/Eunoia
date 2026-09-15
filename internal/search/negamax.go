package search

import (
	"math"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/tt"
)

const (
	nmpMinDepth = 3

	lmrMinDepth = 3
	lmrMinMoves = 2

	lmpBase       = 4
	lmpMultiplier = 3
	lmpMaxDepth   = 4
)

func (ss *SearchState) negamax(pos board.Position, depth, ply int, alpha, beta int16) int16 {
	ss.pv.Init(ply)

	if ss.ShouldStop(Hard) {
		return 0
	}

	ss.Nodes++

	isRoot := ply == 0
	isPV := beta > alpha+1

	if !isRoot && ss.isDraw(&pos) {
		return 0
	}

	alphaOrig := alpha

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	if ttHit && !isRoot && !isPV && entry.Depth >= uint8(depth) {
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
		return ss.qsearch(pos, alpha, beta)
	}

	mover := pos.SideToMove
	inCheck := movegen.InCheck(&pos, mover)

	// reverse futility pruning
	staticEval := evaluate(&pos)
	margin := 150 * int16(depth)
	if !isPV && !ttHit && !inCheck && staticEval >= beta+margin {
		return staticEval
	}

	// null move pruning
	if !inCheck && depth >= nmpMinDepth {
		reduction := 3
		newPos := pos.MakeNullMove()
		score := -ss.negamax(newPos, depth-reduction, ply+1, -beta, -beta+1)
		if score >= beta {
			return score
		}
	}

	bestValue := -INF
	var bestMove board.Move
	movesSearched := 0

	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(&pos, &movelist)
	ss.orderMoves(&pos, &movelist)

	var quietsTried movegen.Movelist

	var score int16

	for i := range movelist.Len {
		move := movelist.Moves[i]
		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}

		isCapture := move.IsCapture()

		// late move pruning
		if !isPV && !isRoot && !isCapture && !inCheck && !isMateScore(bestValue) &&
			depth <= lmpMaxDepth && movesSearched >= lmpBase+lmpMultiplier*depth*depth {
			continue
		}

		ss.keyHistory = append(ss.keyHistory, newPos.Hash)
		newDepth := depth - 1

		// late move reductions
		isReduced := false
		if movesSearched >= lmrMinMoves && depth >= lmrMinDepth &&
			!isCapture {
			reduction := int(0.99 + math.Log(float64(newDepth))*math.Log(float64(movesSearched))/3.14)

			if reduction > 0 {
				reducedDepth := max(newDepth-reduction, 1)

				// search "late" moves with a
				// reduced depth, in a null window
				score = -ss.negamax(newPos, reducedDepth, ply+1, -alpha-1, -alpha)

				if score <= alpha {
					// this move is trash, don't search it in PVS
					isReduced = true
				}
			}
		}

		// principal variation search
		if !isReduced {
			if movesSearched == 0 {
				// full window search for principal variation
				score = -ss.negamax(newPos, newDepth, ply+1, -beta, -alpha)
			} else {
				// null window search for non-PV line
				score = -ss.negamax(newPos, newDepth, ply+1, -alpha-1, -alpha)
				if alpha < score && score < beta {
					// null window failed, re-search with full window
					score = -ss.negamax(newPos, newDepth, ply+1, -beta, -alpha)
				}
			}
		}
		ss.keyHistory = ss.keyHistory[:len(ss.keyHistory)-1]

		if ss.ShouldStop(Hard) {
			return 0
		}

		movesSearched++

		if score > bestValue {
			bestValue = score
			bestMove = move
			if score > alpha {
				alpha = score
				ss.pv.Store(ply, move)
			}
		}
		if score >= beta { // beta cutoff
			if !isCapture {
				bonus := 300*depth - 250
				ss.updateButterflyHistory(mover, move.From(), move.To(), bonus)

				for i := range quietsTried.Len {
					// penalize quiets that didn't cause beta cutoff
					quietMove := quietsTried.Moves[i]
					ss.updateButterflyHistory(mover, quietMove.From(), quietMove.To(), -bonus)
				}
			}
			break
		}

		if !isCapture {
			quietsTried.Add(move)
		}
	}

	if movesSearched == 0 {
		if inCheck {
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
