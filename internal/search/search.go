package search

import (
	"math"

	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/movegen"
)

func SearchBestMove(pos board.Position, depth int8) board.Move {
	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(&pos, &movelist)
	mover := pos.SideToMove

	var bestMove board.Move
	bestScore := int16(math.MinInt16)

	for i := range movelist.Len {
		move := movelist.Moves[i]
		newPos := pos.MakeMove(move)

		if movegen.InCheck(&newPos, mover) {
			continue
		}
		score := -Negamax(newPos, depth-1)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}
	return bestMove
}

func Negamax(pos board.Position, depth int8) int16 {
	var movelist movegen.Movelist
	movegen.GeneratePseudolegalMoves(&pos, &movelist)
	mover := pos.SideToMove

	max := int16(math.MinInt16)
	legalMoves := 0

	for i := range movelist.Len {
		move := movelist.Moves[i]
		newPos := pos.MakeMove(move)
		if movegen.InCheck(&newPos, mover) {
			continue
		}

		legalMoves++

		if depth <= 0 {
			// count legal moves without recursing
			continue
		}

		score := -Negamax(newPos, depth-1)
		if score > max {
			max = score
		}
	}

	if legalMoves == 0 {
		if movegen.InCheck(&pos, mover) {
			// checkmate
			return -30000 + int16(pos.Ply)
		}
		// stalemate
		return 0
	}

	if depth <= 0 {
		return evaluate(pos)
	}
	return max
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
