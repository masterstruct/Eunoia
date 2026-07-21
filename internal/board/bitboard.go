package board

type Bitboard uint64

var SquareBB [64]Bitboard

func InitBitboards() {
	for sq := range 64 {
		SquareBB[sq] = Bitboard(1) << sq
	}
}

func (bb *Bitboard) SetBit(sq uint8) {
	*bb |= SquareBB[sq]
}

func (bb *Bitboard) ClearBit(sq uint8) {
	*bb &^= SquareBB[sq]
}

func (bb Bitboard) IsBitSet(sq uint8) bool {
	return (bb & SquareBB[sq]) != 0
}
