package board

type PRNG struct {
	state uint64
}

// Splitmix64 pseudorandom number generator
func (p *PRNG) Next() uint64 {
	p.state += 0x9e3779b97f4a7c15
	z := p.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

type Zobrist struct {
	Piece      [2][6][64]uint64 // [color][pieceType][sq]
	Castling   [16]uint64
	EnPassant  [8]uint64
	SideToMove uint64
}

func NewZobrist(rng *PRNG) *Zobrist {
	z := &Zobrist{}
	for color := range NoColor {
		for pt := range PieceTypes() {
			for sq := range NoSquare {
				z.Piece[color][pt][sq] = rng.Next()
			}
		}
	}
	for i := range 16 {
		z.Castling[i] = rng.Next()
	}
	for i := range 8 {
		z.EnPassant[i] = rng.Next()
	}
	z.SideToMove = rng.Next()
	return z
}
