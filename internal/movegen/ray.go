package movegen

import "github.com/masterstruct/Eunoia/internal/board"

// Thank you Dan, creator of Hobbes, for the ray implementation this file is based on
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/board/ray.rs

var (
	between   [64][64]board.Bitboard
	extending [64][64]board.Bitboard
	beyond    [64][64]board.Bitboard
)

func Between(a, b board.Square) board.Bitboard {
	return between[a][b]
}

func Extending(a, b board.Square) board.Bitboard {
	return extending[a][b]
}

func Beyond(a, b board.Square) board.Bitboard {
	return beyond[a][b]
}

func init() {
	initBetween()
	initExtending()
	initBeyond()
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

func initExtending() {
	for a := range board.NoSquare {
		for b := range board.NoSquare {
			if RookAttacks(a, board.EmptyBB).IsBitSet(b) {
				extending[a][b] = RookAttacks(a, board.EmptyBB) & RookAttacks(b, board.EmptyBB)
			}

			if BishopAttacks(a, board.EmptyBB).IsBitSet(b) {
				extending[a][b] = BishopAttacks(a, board.EmptyBB) & BishopAttacks(b, board.EmptyBB)
			}
		}
	}
}

func initBeyond() {
	for a := range board.NoSquare {
		for b := range board.NoSquare {
			if RookAttacks(a, board.EmptyBB).IsBitSet(b) {
				beyond[a][b] = RookAttacks(a, board.EmptyBB) & RookAttacks(b, board.SquareBB[a]) &^ Between(a, b)
			}

			if BishopAttacks(a, board.EmptyBB).IsBitSet(b) {
				beyond[a][b] = BishopAttacks(a, board.EmptyBB) & BishopAttacks(b, board.SquareBB[a]) &^ Between(a, b)
			}
		}
	}
}
