package movegen

import "github.com/masterstruct/Eunoia/internal/board"

var (
	rookDirs   = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	bishopDirs = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
)

func RookAttacks(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return rayAttacks(sq, rookDirs, occupied)
}

func rayAttacks(sq board.Square, dirs [][2]int, occupied board.Bitboard) board.Bitboard {
	bb := board.EmptyBB

	for _, dir := range dirs {
		df, dr := dir[0], dir[1]
		currSq := sq

		for currSq.IsValid() && !occupied.IsBitSet(currSq) {
			newFile := currSq.File() + df
			newRank := currSq.Rank() + dr

			currSq = board.NewSquare(newFile, newRank)
			if currSq.IsValid() {
				bb.SetBit(currSq)
			}
		}
	}
	return bb
}
