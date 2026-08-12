package movegen

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

func IsSquareAttacked(pos board.Position, sq board.Square, byColor board.Color) bool {
	if KnightAttacks[sq]&pos.PieceBB(board.NewPiece(board.Knight, byColor)) != 0 {
		return true
	}
	if KingAttacks[sq]&pos.PieceBB(board.NewPiece(board.King, byColor)) != 0 {
		return true
	}
	if PawnAttacks[byColor.Opponent()][sq]&pos.PieceBB(board.NewPiece(board.Pawn, byColor)) != 0 {
		return true
	}
	if RookAttacks(sq, pos.Occupied())&pos.PieceBB(board.NewPiece(board.Rook, byColor)) != 0 {
		return true
	}
	if BishopAttacks(sq, pos.Occupied())&pos.PieceBB(board.NewPiece(board.Bishop, byColor)) != 0 {
		return true
	}
	if QueenAttacks(sq, pos.Occupied())&pos.PieceBB(board.NewPiece(board.Queen, byColor)) != 0 {
		return true
	}
	return false
}
