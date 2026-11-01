package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/tt"
)

const (
	nmpMinDepth = 3

	lmpBase       = 4
	lmpMultiplier = 3
	lmpMaxDepth   = 4

	fpMaxDepth = 3

	hpMaxDepth   = 4
	hpMultiplier = -2048

	rfpMultiplier    = 150
	rfpImpMultiplier = 50
)

func (ss *SearchState) negamax(pos board.Position, depth int, ply uint16, alpha, beta int32) int32 {
	ss.pv.Clear(ply)
	ss.selDepth = max(ss.selDepth, ply)

	if ss.ShouldStop(Hard) {
		return 0
	}

	isRoot := ply == 0
	isPV := beta > alpha+1

	if !isRoot && ss.isDraw(&pos) {
		return 0
	}

	if ply >= MaxPly {
		return evaluate(&pos)
	}

	// mate distance pruning
	alpha = max(alpha, matedIn(ply))
	beta = min(beta, mateIn(ply))
	if alpha >= beta {
		return alpha
	}

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	if ttHit && !isPV && entry.Depth >= uint8(depth) {
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
	ttMove := board.NullMove
	if ttHit {
		ttMove = entry.Move
	}

	stm := pos.SideToMove
	inCheck := pos.InCheck()

	if depth <= 0 && !inCheck {
		return ss.qsearch(pos, ply, alpha, beta)
	}

	depth = max(depth, 0)

	// static eval
	rawEval := -INF
	staticEval := -INF
	correction := int32(0)
	if !inCheck {
		rawEval = evaluate(&pos)
		correction = ss.corrHist.Correction(&pos)
		staticEval = rawEval + correction
	}
	ss.stack[ply].staticEval = staticEval

	improving := ss.improvement(ply, staticEval, inCheck) > 0
	improvingInt := 0
	if improving {
		improvingInt = 1
	}

	// reverse futility pruning
	rfpMargin := rfpMargin(depth, improving)
	if !isPV && !ttHit && !inCheck && staticEval >= beta+rfpMargin {
		return staticEval
	}

	// null move pruning
	if !inCheck && !isPV && pos.HasNonPawns() && staticEval >= beta && depth >= nmpMinDepth {
		reduction := 3 + depth/5
		reducedDepth := depth - reduction

		newPos := pos
		newPos.MakeNullMove()
		ss.tt.Prefetch(newPos.Hash)
		ss.Nodes++

		ss.Keys.Push(newPos.Hash)
		score := -ss.negamax(newPos, reducedDepth, ply+1, -beta, -beta+1)
		ss.Keys.Pop()

		if ss.ShouldStop(Hard) {
			return 0
		}

		if score >= beta {
			// at low depths we can return NMP result directly
			if depth <= 14 || ss.nmpMinPly > 0 {
				if isMateScore(score) {
					return beta
				} else {
					return score
				}
			}

			// verification search
			ss.nmpMinPly = ply + 3*uint16(reducedDepth)/4
			verifScore := ss.negamax(pos, reducedDepth, ply, beta-1, beta)
			ss.nmpMinPly = 0

			if ss.ShouldStop(Hard) {
				return 0
			}

			if verifScore >= beta {
				return score
			}
		}
	}

	bestScore := -INF
	bestMove := board.NullMove
	legalMoves := 0
	searchedMoves := 0
	flag := tt.Upper

	movePicker := NewMovePicker(ply, ttMove, false)

	var quietsTried board.Movelist

	var score int32

	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		legalMoves++

		isQuiet := !move.IsNoisy()
		isMated := isMated(bestScore)
		isKiller := ss.killers[ply] == move

		var historyScore int32
		if isQuiet {
			historyScore = ss.butterflyHistory[stm][move.From()][move.To()]
		}

		// futility pruning
		if !isRoot && isQuiet && !inCheck && !isMated &&
			depth < fpMaxDepth && staticEval+150 <= alpha {
			movePicker.skipQuiets = true
			continue
		}

		// late move pruning
		if !isPV && isQuiet && !inCheck && !isMated &&
			depth <= lmpMaxDepth && searchedMoves >= lmpBase+(lmpMultiplier*depth*depth)/(2-improvingInt) {
			movePicker.skipQuiets = true
			continue
		}

		// history pruning
		if !isRoot && !isMated && !isKiller && isQuiet && depth <= hpMaxDepth &&
			historyScore < hpMultiplier*int32(depth*depth) {
			movePicker.skipQuiets = true
			continue
		}

		// SEE pruning
		var seeThreshold int32
		if isQuiet {
			seeThreshold = seeQuietMultiplier * int32(depth)
		} else {
			seeThreshold = seeNoisyMultiplier * int32(depth)
		}
		if !isMated && !see(&pos, move, seeThreshold) {
			continue
		}

		newPos := pos
		ss.makeMove(&newPos, move)

		newDepth := depth - 1

		// late move reductions
		isReduced := false
		if searchedMoves >= lmrMinMoves && depth >= lmrMinDepth &&
			isQuiet {
			reduction := lmr[min(newDepth, lmrMaxDepth)][min(searchedMoves, lmrMaxMoves)]
			if !improving {
				reduction += 1
			}

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
			if searchedMoves == 0 {
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

		ss.unmakeMove()
		searchedMoves++

		if ss.ShouldStop(Hard) {
			return 0
		}

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if score > alpha {
			alpha = score
			flag = tt.Exact
			if isPV {
				ss.pv.Store(ply, move)
			}
		}

		if score >= beta { // beta cutoff
			flag = tt.Lower
			if isQuiet {
				ss.killers[ply] = bestMove

				bonus := 300*int32(depth) - 250
				ss.updateButterflyHistory(stm, move.From(), move.To(), bonus)

				for i := range quietsTried.Len {
					// penalize quiets that didn't cause beta cutoff
					quietMove := quietsTried.Moves[i].Move
					ss.updateButterflyHistory(stm, quietMove.From(), quietMove.To(), -bonus)
				}
			}
			break
		}

		if isQuiet {
			quietsTried.Add(move)
		}
	}

	if searchedMoves == 0 {
		if legalMoves > 0 {
			return alpha
		}
		if inCheck {
			// checkmate
			return matedIn(ply)
		}
		// stalemate
		return 0
	}

	// update corrhist
	if !inCheck && !bestMove.IsNoisy() && !isMateScore(bestScore) &&
		!(flag == tt.Upper && bestScore >= staticEval) &&
		!(flag == tt.Lower && bestScore <= staticEval) {
		ss.corrHist.Update(&pos, depth, staticEval, bestScore)
	}

	ss.tt.Store(pos.Hash, bestMove, scoreToTT(bestScore, ply), uint8(depth), flag)

	return bestScore
}

func (ss *SearchState) improvement(ply uint16, staticEval int32, inCheck bool) int32 {
	if !inCheck && ply >= 2 && isDefined(ss.stack[ply-2].staticEval) {
		return staticEval - ss.stack[ply-2].staticEval
	}
	if !inCheck && ply >= 4 && isDefined(ss.stack[ply-4].staticEval) {
		return staticEval - ss.stack[ply-4].staticEval
	}
	return 0
}

func rfpMargin(depth int, improving bool) int32 {
	if improving {
		return int32(depth) * rfpImpMultiplier
	}
	return int32(depth) * rfpMultiplier
}
