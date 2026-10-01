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

func rookAttacksSlow(sq Square, occ Bitboard) Bitboard {
	return rayAttacks(sq, rookDirs, occ)
}

func bishopAttacksSlow(sq Square, occ Bitboard) Bitboard {
	return rayAttacks(sq, bishopDirs, occ)
}

func rayAttacks(sq Square, dirs [][2]int, occ Bitboard) Bitboard {
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
			if occ.IsBitSet(currSq) {
				break
			}
		}
	}
	return bb
}

func RookAttacks(sq Square, occ Bitboard) Bitboard {
	return RookMoves[MagicIndex(&RookMagics[sq], occ)]
}

func BishopAttacks(sq Square, occ Bitboard) Bitboard {
	return BishopMoves[MagicIndex(&BishopMagics[sq], occ)]
}

func QueenAttacks(sq Square, occ Bitboard) Bitboard {
	return RookAttacks(sq, occ) | BishopAttacks(sq, occ)
}
