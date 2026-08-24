package search

import (
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
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
		return evaluate(pos)
	}
	ss.Nodes++

	if depth <= 0 {
		return evaluate(pos)
	}

	bestValue := -INF
	legalMoves := 0

	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(&pos, &movelist)
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

	return bestValue
}

func evaluate(pos board.Position) int16 {
	// piece count
	blackPieceCount := pos.Bitboards.Colors[board.Black].CountBits()
	whitePieceCount := pos.Bitboards.Colors[board.White].CountBits()
	score := 50 * (whitePieceCount - blackPieceCount)

	// material count
	score += 100 * (pos.PieceBB(board.WhitePawn).CountBits() - pos.PieceBB(board.BlackPawn).CountBits())
	score += 300 * (pos.PieceBB(board.WhiteKnight).CountBits() - pos.PieceBB(board.BlackKnight).CountBits())
	score += 300 * (pos.PieceBB(board.WhiteBishop).CountBits() - pos.PieceBB(board.BlackBishop).CountBits())
	score += 500 * (pos.PieceBB(board.WhiteRook).CountBits() - pos.PieceBB(board.BlackRook).CountBits())
	score += 900 * (pos.PieceBB(board.WhiteQueen).CountBits() - pos.PieceBB(board.BlackQueen).CountBits())

	// mobility
	var whiteMoves, blackMoves movegen.Movelist

	wPos := pos
	wPos.SideToMove = board.White
	movegen.GeneratePseudolegalMoves(&wPos, &whiteMoves)

	bPos := pos
	bPos.SideToMove = board.Black
	movegen.GeneratePseudolegalMoves(&bPos, &blackMoves)

	score += 5 * (whiteMoves.Len - blackMoves.Len)

	// giving checks is good
	if movegen.InCheck(&pos, board.Black) {
		score += 50
	}
	if movegen.InCheck(&pos, board.White) {
		score -= 50
	}

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
