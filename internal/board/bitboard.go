package board

import "math/bits"

type Bitboard uint64

const FullBB Bitboard = ^Bitboard(0)
const EmptyBB Bitboard = 0

var SquareBB [64]Bitboard

func init() {
	for sq := range 64 {
		SquareBB[sq] = Bitboard(1) << sq
	}
}

func (bb *Bitboard) SetBit(sq Square) {
	*bb |= SquareBB[sq]
}

func (bb *Bitboard) ClearBit(sq Square) {
	*bb &^= SquareBB[sq]
}

func (bb Bitboard) IsBitSet(sq Square) bool {
	return (bb & SquareBB[sq]) != 0
}

func (bb Bitboard) CountBits() int {
	return bits.OnesCount64(uint64(bb))
}

func (bb Bitboard) LSB() Square {
	if bb == EmptyBB {
		return NoSquare
	}
	return Square(bits.TrailingZeros64(uint64(bb)))
}

func (bb *Bitboard) PopLSB() Square {
	bit := bb.LSB()
	if bit >= 0 {
		bb.ClearBit(bit)
	}
	return bit
}
