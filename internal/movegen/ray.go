package movegen

import "github.com/masterstruct/Eunoia/internal/board"

var (
	between [64][64]board.Bitboard
)

func Between(a, b board.Square) board.Bitboard {
	return between[a][b]
}

func init() {
	initBetween()
}

func initBetween() {
	for a := range board.NoSquare {
		for b := range board.NoSquare {
			if RookAttacks(a, board.EmptyBB).IsBitSet(b) {
				between[a][b] = RookAttacks(a, board.SquareBB[b]) & RookAttacks(b, board.SquareBB[a])
			}

			if BishopAttacks(a, board.EmptyBB).IsBitSet(b) {
				between[a][b] = BishopAttacks(a, board.SquareBB[b]) & BishopAttacks(b, board.SquareBB[a])
			}
		}
	}
}
