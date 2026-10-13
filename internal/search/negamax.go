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
)

func (ss *SearchState) negamax(pos board.Position, depth int, ply uint16, alpha, beta int32) int32 {
	ss.pv.Clear(ply)

	if ss.ShouldStop(Hard) {
		return 0
	}

	isRoot := ply == 0
	isPV := beta > alpha+1

	if !isRoot && ss.isDraw(&pos) {
		return 0
	}

	// mate distance pruning
	alpha = max(alpha, matedIn(ply))
	beta = min(beta, mateIn(ply))
	if alpha >= beta {
		return alpha
	}

	alphaOrig := alpha

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

	if depth <= 0 {
		return ss.qsearch(pos, ply, alpha, beta)
	}

	mover := pos.SideToMove
	inCheck := pos.InCheck()

	// reverse futility pruning
	staticEval := evaluate(&pos)
	margin := 150 * int32(depth)
	if !isPV && !ttHit && !inCheck && staticEval >= beta+margin {
		return staticEval
	}

	// null move pruning
	if !inCheck && !isPV && staticEval >= beta && depth >= nmpMinDepth {
		reduction := 3

		newPos := pos.MakeNullMove()
		ss.Nodes++

		ss.Keys.Push(newPos.Hash)
		score := -ss.negamax(newPos, depth-reduction, ply+1, -beta, -beta+1)
		ss.Keys.Pop()

		if ss.ShouldStop(Hard) {
			return 0
		}

		if score >= beta {
			if isMateScore(score) {
				score = beta
			}
			return score
		}
	}

	bestScore := -INF
	var bestMove board.Move

	legalMoves := 0
	searchedMoves := 0

	movePicker := NewMovePicker(ttMove, false)

	var quietsTried board.Movelist

	var score int32

	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		legalMoves++

		isQuiet := !move.IsNoisy()
		isMated := isMated(bestScore)

		// futility pruning
		if !isRoot && isQuiet && !inCheck && !isMated &&
			depth < fpMaxDepth && staticEval+150 <= alpha {
			movePicker.skipQuiets = true
			continue
		}

		// late move pruning
		if !isPV && isQuiet && !inCheck && !isMated &&
			depth <= lmpMaxDepth && searchedMoves >= lmpBase+lmpMultiplier*depth*depth {
			movePicker.skipQuiets = true
			continue
		}

		newPos := pos.MakeMove(move)
		ss.Keys.Push(newPos.Hash)
		ss.Nodes++

		newDepth := depth - 1

		// late move reductions
		isReduced := false
		if searchedMoves >= lmrMinMoves && depth >= lmrMinDepth &&
			isQuiet {
			reduction := lmr[min(newDepth, lmrMaxDepth)][min(searchedMoves, lmrMaxMoves)]

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

		ss.Keys.Pop()
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
			if isPV {
				ss.pv.Store(ply, move)
			}
		}

		if score >= beta { // beta cutoff
			if isQuiet {
				bonus := 300*int32(depth) - 250
				ss.updateButterflyHistory(mover, move.From(), move.To(), bonus)

				for i := range quietsTried.Len {
					// penalize quiets that didn't cause beta cutoff
					quietMove := quietsTried.Moves[i].Move
					ss.updateButterflyHistory(mover, quietMove.From(), quietMove.To(), -bonus)
				}
			}
			break
		}

		if isQuiet {
			quietsTried.Add(move)
		}
	}

	if legalMoves == 0 {
		if inCheck {
			// checkmate
			return matedIn(ply)
		}
		// stalemate
		return 0
	}

	flag := tt.FlagFromScore(bestScore, alphaOrig, beta)
	ss.tt.Store(pos.Hash, bestMove, scoreToTT(bestScore, ply), uint8(depth), flag)

	return bestScore
}
