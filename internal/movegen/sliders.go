package movegen

import "github.com/masterstruct/Eunoia/internal/board"

var (
	rookDirs   = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	bishopDirs = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
)

func RookAttacks(sq board.Square, occupied board.Bitboard) board.Bitboard {
	return board.EmptyBB
}
