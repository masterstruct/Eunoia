package search

import (
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
		bestScore := -INF

		var movelist movegen.Movelist
		movegen.GeneratePseudolegalMoves(&pos, &movelist)
		ss.OrderMoves(&pos, &movelist)
		mover := pos.SideToMove

		for i := range movelist.Len {
			if ss.searchStopped() {
				return bestMove
			}

			move := movelist.Moves[i]
			newPos := pos.MakeMove(move)

			if movegen.InCheck(&newPos, mover) {
				continue
			}

			score := -ss.Negamax(newPos, depth-1, -INF, INF)

			if score > bestScore {
				bestScore = score
				bestMove = move
			}
		}
	}
	return bestMove
}

func (ss *SearchState) Negamax(pos board.Position, depth int, alpha, beta int16) int16 {
	if ss.searchStopped() {
		return ss.qsearch(&pos, alpha, beta)
	}
	ss.Nodes++

	alphaOrig := alpha

	// TT lookup
	entry, ttHit := ss.tt.Probe(pos.Hash)
	if ttHit && entry.Depth >= uint8(depth) {
		switch entry.Flag {
		case tt.Exact:
			return entry.Score
		case tt.Lower:
			if entry.Score >= beta {
				return entry.Score
			}
		case tt.Upper:
			if entry.Score <= alpha {
				return entry.Score
			}
		}
	}

	if depth <= 0 {
		return evaluate(&pos)
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

		score := -ss.Negamax(newPos, depth-1, -beta, -alpha)
		if score > bestValue {
			bestValue = score
			bestMove = move
			if score > alpha {
				alpha = score
			}
		}
		if score >= beta {
			break
		}
	}

	if legalMoves == 0 {
		if movegen.InCheck(&pos, mover) {
			// checkmate
			return -MATE + int16(pos.Ply)
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

	ss.tt.Store(pos.Hash, bestMove, bestValue, uint8(depth), flag)

	return bestValue
}

func (ss *SearchState) qsearch(pos *board.Position, alpha, beta int16) int16 {
	standPat := evaluate(pos)

	if standPat >= beta {
		return standPat
	}
	if standPat > alpha {
		alpha = standPat
	}

	ss.Nodes++

	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(pos, &movelist)
	ss.OrderMoves(pos, &movelist)
	mover := pos.SideToMove

	for i := range movelist.Len {
		move := movelist.Moves[i]
		if !move.IsCapture() {
			continue
		}

		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}

		score := -ss.qsearch(&newPos, -beta, -alpha)

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
	// piece count
	blackPieceCount := pos.Bitboards.Colors[board.Black].CountBits()
	whitePieceCount := pos.Bitboards.Colors[board.White].CountBits()
	score := 50 * (whitePieceCount - blackPieceCount)

	// mobility
	var whiteMoves, blackMoves movegen.Movelist

	wPos := *pos
	wPos.SideToMove = board.White
	movegen.GeneratePseudolegalMoves(&wPos, &whiteMoves)

	bPos := *pos
	bPos.SideToMove = board.Black
	movegen.GeneratePseudolegalMoves(&bPos, &blackMoves)

	score += 5 * (whiteMoves.Len - blackMoves.Len)

	// giving checks is good
	if movegen.InCheck(pos, board.Black) {
		score += 50
	}
	if movegen.InCheck(pos, board.White) {
		score -= 50
	}

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
