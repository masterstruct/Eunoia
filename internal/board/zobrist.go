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
	piece      [2][6][64]uint64 // [color][pieceType][sq]
	castling   [16]uint64
	enPassant  [8]uint64
	sideToMove uint64
}

func NewZobrist(rng *PRNG) *Zobrist {
	z := &Zobrist{}
	for color := range NoColor {
		for pt := range PieceTypes() {
			for sq := range NoSquare {
				z.piece[color][pt][sq] = rng.Next()
			}
		}
	}
	for i := range 16 {
		z.castling[i] = rng.Next()
	}
	for i := range 8 {
		z.enPassant[i] = rng.Next()
	}
	z.sideToMove = rng.Next()
	return z
}

func (z *Zobrist) PieceKey(color Color, piece PieceType, square Square) uint64 {
	return z.piece[color][piece][square]
}

func (z *Zobrist) CastlingKey(cr CastlingRights) uint64 {
	return z.castling[cr]
}

func (z *Zobrist) EnPassantKey(file File) uint64 {
	return z.enPassant[file]
}

func (z *Zobrist) SideToMoveKey() uint64 {
	return z.sideToMove
}
