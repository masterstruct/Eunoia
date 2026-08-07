package movegen

import "github.com/masterstruct/Eunoia/internal/board"

var (
	KnightAttacks [64]board.Bitboard
	KingAttacks   [64]board.Bitboard
	PawnAttacks   [2][64]board.Bitboard // indexed by board.Color
)
