package search

import "github.com/masterstruct/Eunoia/internal/board"

const (
	seePawnValue   = 100
	seeKnightValue = 300
	seeBishopValue = 330
	seeRookValue   = 500
	seeQueenValue  = 900
)

func seeValue(pt board.PieceType) int32 {
	switch pt {
	case board.Pawn:
		return seePawnValue
	case board.Knight:
		return seeKnightValue
	case board.Bishop:
		return seeBishopValue
	case board.Rook:
		return seeRookValue
	case board.Queen:
		return seeQueenValue
	default:
		return 0
	}
}

func seeMoveValue(pos *board.Position, move board.Move) int32 {
	value := seeValue(pos.Board[move.To()].Type)
	if move.IsPromo() {
		value += seeValue(move.Promo())
	} else if move.IsEnPassant() {
		value = seeValue(board.Pawn)
	}
	return value
}

func seeSquareAttackers(pos *board.Position, sq board.Square, occ board.Bitboard) board.Bitboard {
	wPawns := board.PawnAttacks[board.Black][sq] & pos.PieceBB(board.WhitePawn)
	bPawns := board.PawnAttacks[board.White][sq] & pos.PieceBB(board.BlackPawn)
	knights := board.KnightAttacks[sq] & pos.Pieces[board.Knight]
	kings := board.KingAttacks[sq] & pos.Pieces[board.King]

	queens := pos.Pieces[board.Queen]
	bishops := pos.Pieces[board.Bishop]
	rooks := pos.Pieces[board.Rook]

	return wPawns | bPawns | knights | kings |
		board.BishopAttacks(sq, occ)&(bishops|queens) |
		board.RookAttacks(sq, occ)&(rooks|queens)
}
