package search

import (
	"os"
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
	"github.com/masterstruct/Eunoia/internal/tt"
)

const (
	MATE int16 = 30000
	INF  int16 = 32000
)

func (ss *SearchState) SearchBestMove(pos board.Position, maxDepth int) board.Move {
	var bestMove board.Move

	// iterative deepening
	for depth := 1; depth <= maxDepth; depth++ {
		if (ss.SoftNodes > 0 && ss.Nodes >= ss.SoftNodes) || (!ss.SoftTime.IsZero() && time.Now().After(ss.SoftTime)) {
			break
		}

		score := ss.Negamax(pos, depth, 0, -INF, INF)

		if len(ss.pv.Line()) == 0 {
			// interruped before first move search completed,
			// use bestMove from previous depth
			break
		}

		bestMove = ss.pv.Line()[0]
		ss.printPV(os.Stdout, depth, score)

		if ss.searchStopped() {
			break
		}
	}
	return bestMove
}

func (ss *SearchState) Negamax(pos board.Position, depth, ply int, alpha, beta int16) int16 {
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
	ss.OrderMoves(&pos, &movelist)
	mover := pos.SideToMove

	for i := range movelist.Len {
		move := movelist.Moves[i]
		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}

		legalMoves++

		score := -ss.Negamax(newPos, depth-1, ply+1, -beta, -alpha)

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

func (ss *SearchState) qsearch(pos *board.Position, alpha, beta int16) int16 {
	ss.Nodes++

	standPat := evaluate(pos)

	if standPat >= beta {
		return standPat
	}
	if standPat > alpha {
		alpha = standPat
	}

	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(pos, &movelist)
	ss.OrderMoves(pos, &movelist)
	mover := pos.SideToMove

	for i := range movelist.Len {
		move := movelist.Moves[i]
		if !move.IsCapture() && !move.IsPromo() {
			continue
		}

		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}

		score := -ss.qsearch(&newPos, -beta, -alpha)

		if ss.searchStopped() {
			return 0
		}

		if score >= beta {
			return score
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}

func evaluate(pos *board.Position) int16 {
	var score int

	// PSQT
	score += evaluatePSQT(pos)

	if pos.SideToMove == board.Black {
		return int16(-score)
	}
	return int16(score)
}

func (ss *SearchState) searchStopped() bool {
	if ss.Stop {
		return true
	}
	if ss.MaxNodes > 0 && ss.Nodes >= ss.MaxNodes {
		return true
	}
	if !ss.MaxTime.IsZero() && time.Now().After(ss.MaxTime) {
		return true
	}
	return false
}

func (ss *SearchState) OrderMoves(pos *board.Position, movelist *movegen.Movelist) {
	n := movelist.Len
	if n == 0 {
		return
	}

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	var ttMove board.Move
	if ttHit {
		ttMove = entry.Move
	}

	// score moves
	// TODO: make separate function
	var scores [movegen.MaxMoves]int16
	for i := range n {
		move := movelist.Moves[i]

		if ttHit && move == ttMove {
			scores[i] = 32000
			continue
		}

		if move.IsCapture() {
			attacker, _ := pos.PieceOn(move.From())
			victim, _ := pos.PieceOn(move.To())

			victimType := victim.Type
			if move.IsEnPassant() {
				victimType = board.Pawn
			}

			scores[i] = mvvlvaScore(victimType, attacker.Type)
		}
	}

	// reverse insertion sort
	for i := 1; i < n; i++ {
		score := scores[i]
		move := movelist.Moves[i]

		j := i - 1
		for j >= 0 && scores[j] < score {
			scores[j+1] = scores[j]
			movelist.Moves[j+1] = movelist.Moves[j]
			j--
		}

		scores[j+1] = score
		movelist.Moves[j+1] = move
	}
}

// formula from https://asteri.sm/files/2023-02-20-viri-wiki#mvvlva
func mvvlvaScore(victim, attacker board.PieceType) int16 {
	return int16(victim)*1000 + 60 - int16(attacker)*10
}
