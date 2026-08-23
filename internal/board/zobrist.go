package board

// just a random value
const zobristSeed uint64 = 16812629825545983742

var ZobristTable = NewZobrist(NewPRNG(zobristSeed))

type PRNG struct {
	state uint64
}

func NewPRNG(seed uint64) *PRNG {
	return &PRNG{state: seed}
}

// Splitmix64 pseudorandom number generator
func (p *PRNG) Next() uint64 {
	p.state += 0x9e3779b97f4a7c15
	z := p.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

type zobrist struct {
	piece      [2][6][64]uint64 // [color][pieceType][sq]
	castling   [16]uint64
	enPassant  [8]uint64
	sideToMove uint64
}

func NewZobrist(rng *PRNG) *zobrist {
	z := &zobrist{}
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

func (z *zobrist) PieceKey(color Color, piece PieceType, square Square) uint64 {
	return z.piece[color][piece][square]
}

func (z *zobrist) CastlingKey(cr CastlingRights) uint64 {
	return z.castling[cr]
}

func (z *zobrist) EnPassantKey(file File) uint64 {
	return z.enPassant[file]
}

func (z *zobrist) SideToMoveKey() uint64 {
	return z.sideToMove
}

func (z *zobrist) ComputeHash(pos *Position) uint64 {
	var hash uint64

	occupied := pos.Occupied()
	for occupied != 0 {
		sq := occupied.PopLSB()

		piece, ok := pos.PieceOn(sq)
		if !ok {
			// looping over occupied, thus
			// the piece cannot be NoPiece.
			// panics if pos.Occupied() and
			// pos.Board are out of sync.
			panic("unreachable")
		}

		hash ^= z.PieceKey(piece.Color, piece.Type, sq)
	}

	hash ^= z.CastlingKey(pos.CastlingRights)
	epSq := pos.EnPassant
	if epSq != NoSquare {
		hash ^= z.EnPassantKey(epSq.File())
	}

	hash ^= z.SideToMoveKey()

	return hash
}
