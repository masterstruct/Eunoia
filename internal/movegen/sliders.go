package movegen

import "github.com/masterstruct/Eunoia/internal/board"

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

func RookAttacksSlow(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return rayAttacks(sq, rookDirs, occupied)
}

func BishopAttacksSlow(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return rayAttacks(sq, bishopDirs, occupied)
}

func QueenAttacksSlow(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return rayAttacks(sq, rookDirs, occupied) | rayAttacks(sq, bishopDirs, occupied)
}

func rayAttacks(sq board.Square, dirs [][2]int, occupied board.Bitboard) board.Bitboard {
	bb := board.EmptyBB

	for _, dir := range dirs {
		df, dr := board.File(dir[0]), board.Rank(dir[1])
		currSq := sq

		for {
			newFile := currSq.File() + df
			newRank := currSq.Rank() + dr

			currSq = board.NewSquare(newFile, newRank)
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

func RookAttacks(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return RookMoves[sq][MagicIndex(&RookMagics[sq], occupied)]
}

func BishopAttacks(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return BishopMoves[sq][MagicIndex(&BishopMagics[sq], occupied)]
}

func QueenAttacks(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return RookAttacks(sq, occupied) | BishopAttacks(sq, occupied)
}
