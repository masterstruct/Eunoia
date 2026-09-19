package board

// Thank you Dan, creator of Hobbes, for the ray implementation this file is based on
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/board/ray.rs

var (
	between   [64][64]Bitboard
	extending [64][64]Bitboard
	beyond    [64][64]Bitboard
)

func Between(a, b Square) Bitboard {
	return between[a][b]
}

func Extending(a, b Square) Bitboard {
	return extending[a][b]
}

func Beyond(a, b Square) Bitboard {
	return beyond[a][b]
}

func init() {
	initBetween()
	initExtending()
	initBeyond()
}

func initBetween() {
	for a := range NoSquare {
		for b := range NoSquare {
			if RookAttacks(a, EmptyBB).IsBitSet(b) {
				between[a][b] = RookAttacks(a, SquareBB[b]) & RookAttacks(b, SquareBB[a])
			}

			if BishopAttacks(a, EmptyBB).IsBitSet(b) {
				between[a][b] = BishopAttacks(a, SquareBB[b]) & BishopAttacks(b, SquareBB[a])
			}
		}
	}
}

func initExtending() {
	for a := range NoSquare {
		for b := range NoSquare {
			if RookAttacks(a, EmptyBB).IsBitSet(b) {
				extending[a][b] = RookAttacks(a, EmptyBB) & RookAttacks(b, EmptyBB)
			}

			if BishopAttacks(a, EmptyBB).IsBitSet(b) {
				extending[a][b] = BishopAttacks(a, EmptyBB) & BishopAttacks(b, EmptyBB)
			}
		}
	}
}

func initBeyond() {
	for a := range NoSquare {
		for b := range NoSquare {
			if RookAttacks(a, EmptyBB).IsBitSet(b) {
				beyond[a][b] = RookAttacks(a, EmptyBB) & RookAttacks(b, SquareBB[a]) &^ Between(a, b)
			}

			if BishopAttacks(a, EmptyBB).IsBitSet(b) {
				beyond[a][b] = BishopAttacks(a, EmptyBB) & BishopAttacks(b, SquareBB[a]) &^ Between(a, b)
			}
		}
	}
}
