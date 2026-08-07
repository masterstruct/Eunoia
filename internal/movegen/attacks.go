package movegen

import "github.com/masterstruct/Eunoia/internal/board"

var (
	KnightAttacks [64]board.Bitboard
	KingAttacks   [64]board.Bitboard
	PawnAttacks   [2][64]board.Bitboard // indexed by board.Color
)

var knightOffsets = [8][2]int{
	{1, 2}, {2, 1}, {2, -1}, {1, -2},
	{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
}

var kingOffsets = [8][2]int{
	{1, 1}, {1, 0}, {1, -1},
	{0, -1}, {-1, -1}, {-1, 0},
	{-1, 1}, {0, 1},
}

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
	return attacksFromOffsets(sq, knightOffsets[:])
}

func kingAttacksFrom(sq board.Square) board.Bitboard {
	return attacksFromOffsets(sq, kingOffsets[:])
}

func pawnAttacksFrom(sq board.Square, c board.Color) board.Bitboard {
	return board.EmptyBB
}

func attacksFromOffsets(sq board.Square, offsets [][2]int) board.Bitboard {
	bb := board.EmptyBB
	for _, offset := range offsets {
		newFile := sq.File() + offset[0]
		newRank := sq.Rank() + offset[1]

		targetSq := board.NewSquare(newFile, newRank)
		if targetSq.IsValid() {
			bb.SetBit(targetSq)
		}
	}
	return bb
}
