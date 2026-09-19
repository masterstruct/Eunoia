package board

var (
	KnightAttacks [64]Bitboard
	KingAttacks   [64]Bitboard
	PawnAttacks   [2][64]Bitboard // indexed by Color
)

var knightOffsets = [8][2]int{
	/**/ {-1, 2}, {1, 2},
	{-2, 1} /*    */, {2, 1},
	/*         ♞          */
	{-2, -1} /*   */, {2, -1},
	/**/ {-1, -2}, {1, -2},
}

var kingOffsets = [8][2]int{
	{-1, 1}, {0, 1}, {1, 1},
	{-1, 0} /*  ♚ */, {1, 0},
	{-1, -1}, {0, -1}, {1, -1},
}

var pawnOffsets = [2][2][2]int{
	{{-1, -1}, {1, -1}}, // Black
	{{-1, 1}, {1, 1}},   // White
}

func init() {
	InitAttackTables()
}

func InitAttackTables() {
	InitBitboards()
	for sq := A1; sq <= H8; sq++ {
		KnightAttacks[sq] = knightAttacksFrom(sq)
		KingAttacks[sq] = kingAttacksFrom(sq)
		PawnAttacks[Black][sq] = pawnAttacksFrom(sq, Black)
		PawnAttacks[White][sq] = pawnAttacksFrom(sq, White)
	}
}

func knightAttacksFrom(sq Square) Bitboard {
	return attacksFromOffsets(sq, knightOffsets[:])
}

func kingAttacksFrom(sq Square) Bitboard {
	return attacksFromOffsets(sq, kingOffsets[:])
}

func pawnAttacksFrom(sq Square, color Color) Bitboard {
	return attacksFromOffsets(sq, pawnOffsets[color][:])
}

func attacksFromOffsets(sq Square, offsets [][2]int) Bitboard {
	bb := EmptyBB
	for _, offset := range offsets {
		df, dr := File(offset[0]), Rank(offset[1])

		newFile := File(sq.File() + df)
		newRank := Rank(sq.Rank() + dr)

		targetSq := NewSquare(newFile, newRank)
		if targetSq.IsValid() {
			bb.SetBit(targetSq)
		}
	}
	return bb
}

func pawnAttacks(pawns Bitboard, color Color) Bitboard {
	if color == White {
		return pawns.northEast() | pawns.northWest()
	}
	return pawns.southEast() | pawns.southWest()
}
