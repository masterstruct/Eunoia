package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/tt"
)

const (
	nmpMinDepth = 3

	lmpBase     = 4
	lmpScale    = 3
	lmpMaxDepth = 4

	fpMaxDepth = 5

	hpMaxDepth = 4
	hpScale    = -2048

	rfpScale    = 150
	rfpImpScale = 50

	fpBase  = 100
	fpScale = 50
)

func (ss *SearchState) negamax(pos board.Position, depth int32, ply uint16, alpha, beta int32) int32 {
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

	singular := ss.stack[ply].excluded
	singularSearch := singular != board.NullMove

	// TT lookup
	ttHit := false
	ttMove := board.NullMove
	ttFlag := tt.Lower
	var ttDepth int32
	ttScore := -INF
	if !singularSearch {
		entry := ss.tt.Probe(pos.Hash)
		ttHit = pos.Hash == entry.Key
		if ttHit {
			ttMove = entry.Move
			ttFlag = entry.Flag
			ttScore = scoreFromTT(entry.Score, ply)
			ttDepth = int32(entry.Depth)
			if !isPV && uint8(ttDepth) >= uint8(depth) { // FIXME: negamax depth can be negative here
				switch ttFlag {
				case tt.Exact:
					return ttScore
				case tt.Lower:
					if ttScore >= beta {
						return ttScore
					}
				case tt.Upper:
					if ttScore <= alpha {
						return ttScore
					}
				}
			}
		}
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
	var correction int32
	if !inCheck {
		rawEval = evaluate(&pos)
		correction = ss.corrHist.Correction(&pos)
		staticEval = clampMate(rawEval + correction)
	}
	ss.stack[ply].staticEval = staticEval

	improving := ss.improvement(ply, staticEval, inCheck) > 0
	var improvingInt32 int32
	if improving {
		improvingInt32 = 1
	}

	// pre-moveloop pruning
	if !isPV && !inCheck && !singularSearch {
		// reverse futility pruning
		rfpMargin := rfpMargin(depth, improving)
		if !ttHit && staticEval >= beta+rfpMargin {
			return staticEval
		}

		// null move pruning
		if pos.HasNonPawns() && staticEval >= beta && depth >= nmpMinDepth {
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
	}

	bestScore := -INF
	bestMove := board.NullMove
	legalMoves := 0
	var searchedMoves int32
	flag := tt.Upper

	contIndices := ss.NewContIndices(ply)
	movePicker := NewMovePicker(ply, ttMove, false, contIndices)

	var quietsTried board.Movelist

	var score int32

	nextMove := func() board.Move { return movePicker.Next(&pos, ss) }
	for move := nextMove(); move != board.NullMove; move = nextMove() {
		legalMoves++
		if singularSearch && move == singular {
			continue
		}

		isQuiet := !move.IsNoisy()
		isMated := isMated(bestScore)
		isKiller := ss.stack[ply].killer == move
		from := move.From()
		to := move.To()

		var extension int32

		var historyScore int32
		if isQuiet {
			historyScore = ss.history.Quiet.Get(stm, from, to)
		}

		// futility pruning
		futilityMargin := fpBase + fpScale*depth
		if !isRoot && isQuiet && !inCheck && !isMated &&
			depth < fpMaxDepth && staticEval+futilityMargin <= alpha {
			movePicker.skipQuiets = true
			continue
		}

		// late move pruning
		if !isPV && isQuiet && !inCheck && !isMated &&
			depth <= lmpMaxDepth && searchedMoves >= lmpBase+(lmpScale*depth*depth)/(2-improvingInt32) {
			movePicker.skipQuiets = true
			continue
		}

		// history pruning
		if !isRoot && !isMated && !isKiller && isQuiet && depth <= hpMaxDepth &&
			historyScore < hpScale*depth*depth {
			movePicker.skipQuiets = true
			continue
		}

		// SEE pruning
		var seeThreshold int32
		if isQuiet {
			seeThreshold = seeQuietScale * depth
		} else {
			seeThreshold = seeNoisyScale * depth
		}
		if !isMated && !see(&pos, move, seeThreshold) {
			continue
		}

		// singular extensions
		if !isRoot && !singularSearch &&
			ttHit && move == ttMove && ttFlag != tt.Upper &&
			depth >= 8 && ttDepth >= depth-3 {

			sBeta := max(ttScore-2*depth, -MATE+1)
			sDepth := (depth - 1) / 2

			ss.stack[ply].excluded = move
			score := ss.negamax(pos, sDepth, ply, sBeta-1, sBeta)
			ss.stack[ply].excluded = board.NullMove

			if score < sBeta {
				extension = 1
			}
		}

		newPos := pos
		ss.makeMove(&newPos, move, ply)

		newDepth := depth + extension - 1

		// principal variation search
		if searchedMoves == 0 {
			score = -ss.negamax(newPos, newDepth, ply+1, -beta, -alpha)
		} else {
			// late move reductions
			var lmr int32
			if depth >= lmrMinDepth && searchedMoves >= lmrMinMoves && isQuiet {
				lmr = getLMR(newDepth, searchedMoves)
				if !improving {
					lmr++
				}
			}

			lmrDepth := min(max(newDepth-lmr, 1), newDepth)
			score = -ss.negamax(newPos, lmrDepth, ply+1, -alpha-1, -alpha)

			if lmr > 0 && score > alpha {
				score = -ss.negamax(newPos, newDepth, ply+1, -alpha-1, -alpha)
			}

			if isPV && score > alpha && score < beta {
				score = -ss.negamax(newPos, newDepth, ply+1, -beta, -alpha)
			}
		}

		ss.unmakeMove(ply)
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
				ss.stack[ply].killer = bestMove
			}

			ss.history.Update(
				&pos,
				contIndices,
				depth,
				bestMove,
				&quietsTried,
			)
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
	if !inCheck && !singularSearch && !bestMove.IsNoisy() && !isMateScore(bestScore) &&
		!(flag == tt.Upper && bestScore >= staticEval) &&
		!(flag == tt.Lower && bestScore <= staticEval) {
		ss.corrHist.Update(&pos, depth, bestScore, staticEval)
	}

	if !singularSearch {
		ss.tt.Store(pos.Hash, bestMove, scoreToTT(bestScore, ply), uint8(depth), flag)
	}

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

func rfpMargin(depth int32, improving bool) int32 {
	if improving {
		return depth * rfpImpScale
	}
	return depth * rfpScale
}
