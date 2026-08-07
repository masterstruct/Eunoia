package movegen

import "github.com/masterstruct/Eunoia/internal/board"

var (
	KnightAttacks [64]board.Bitboard
	KingAttacks   [64]board.Bitboard
	PawnAttacks   [2][64]board.Bitboard // indexed by board.Color
)

func init() {
	InitAttackTables()
}

func InitAttackTables() {
	for sq := board.A1; sq <= board.H8; sq++ {
		KnightAttacks[sq] = knightAttacksFrom(sq)
		KingAttacks[sq] = kingAttacksFrom(sq)
		PawnAttacks[board.Black][sq] = pawnAttacksFrom(sq, board.Black)
		PawnAttacks[board.White][sq] = pawnAttacksFrom(sq, board.White)
	}
}

func knightAttacksFrom(sq board.Square) board.Bitboard {
	return board.EmptyBB
}

func kingAttacksFrom(sq board.Square) board.Bitboard {
	return board.EmptyBB
}

func pawnAttacksFrom(sq board.Square, c board.Color) board.Bitboard {
	return board.EmptyBB
}
