package board

// Thank you Dan, creator of Hobbes, for the ray implementation this file is based on
// https://github.com/kelseyde/hobbes-chess-engine/blob/main/src/board/ray.rs

var (
	between      [64][64]Bitboard
	extending    [64][64]Bitboard
	beyond       [64][64]Bitboard
	diagonalRays [4][64]Bitboard
)

type DiagonalDirection uint8

const (
	NorthWest DiagonalDirection = iota
	NorthEast
	SouthEast
	SouthWest
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

func Diagonal(sq Square, direction DiagonalDirection) Bitboard {
	return diagonalRays[direction][sq]
}

func init() {
	initBetween()
	initExtending()
	initBeyond()
	initDiagonals()
}

func initBetween() {
	for a := range NoSquare {
		for b := range NoSquare {
			if RookAttacks(a, EmptyBB).IsBitSet(b) {
				between[a][b] = RookAttacks(a, b.Bit()) & RookAttacks(b, a.Bit())
			}

			if BishopAttacks(a, EmptyBB).IsBitSet(b) {
				between[a][b] = BishopAttacks(a, b.Bit()) & BishopAttacks(b, a.Bit())
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
				beyond[a][b] = RookAttacks(a, EmptyBB) & RookAttacks(b, a.Bit()) &^ Between(a, b)
			}

			if BishopAttacks(a, EmptyBB).IsBitSet(b) {
				beyond[a][b] = BishopAttacks(a, EmptyBB) & BishopAttacks(b, a.Bit()) &^ Between(a, b)
			}
		}
	}
}

func initDiagonals() {
	directions := [4][2]int{
		{-1, 1},
		{1, 1},
		{1, -1},
		{-1, -1},
	}
	for sq := A1; sq <= H8; sq++ {
		for direction, step := range directions {
			file := int(sq.File()) + step[0]
			rank := int(sq.Rank()) + step[1]
			for {
				next := NewSquare(File(file), Rank(rank))
				if !next.IsValid() {
					break
				}
				diagonalRays[direction][sq].SetBit(next)
				file += step[0]
				rank += step[1]
			}
		}
	}
}
