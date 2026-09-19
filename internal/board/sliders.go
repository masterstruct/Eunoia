package board

var (
	rookDirs = [][2]int{
		/*     */ {0, 1},
		{-1, 0} /*  ♜ */, {1, 0},
		/*     */ {0, -1},
	}
	bishopDirs = [][2]int{
		{-1, 1}, {1, 1},
		/*      ♝     */
		{-1, -1}, {1, -1},
	}
)

func rookAttacksSlow(sq Square, occupied Bitboard) Bitboard {
	return rayAttacks(sq, rookDirs, occupied)
}

func bishopAttacksSlow(sq Square, occupied Bitboard) Bitboard {
	return rayAttacks(sq, bishopDirs, occupied)
}

func rayAttacks(sq Square, dirs [][2]int, occupied Bitboard) Bitboard {
	bb := EmptyBB

	for _, dir := range dirs {
		df, dr := File(dir[0]), Rank(dir[1])
		currSq := sq

		for {
			newFile := currSq.File() + df
			newRank := currSq.Rank() + dr

			currSq = NewSquare(newFile, newRank)
			if !currSq.IsValid() {
				break
			}

			bb.SetBit(currSq)
			if occupied.IsBitSet(currSq) {
				break
			}
		}
	}
	return bb
}

func RookAttacks(sq Square, occupied Bitboard) Bitboard {
	return RookMoves[MagicIndex(&RookMagics[sq], occupied)]
}

func BishopAttacks(sq Square, occupied Bitboard) Bitboard {
	return BishopMoves[MagicIndex(&BishopMagics[sq], occupied)]
}

func QueenAttacks(sq Square, occupied Bitboard) Bitboard {
	return RookAttacks(sq, occupied) | BishopAttacks(sq, occupied)
}
