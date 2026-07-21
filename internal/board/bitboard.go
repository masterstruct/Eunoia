package board

import "math/bits"

type Bitboard uint64

const FullBB Bitboard = ^Bitboard(0)
const EmptyBB Bitboard = 0

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

func (bb Bitboard) CountBits() int {
	return bits.OnesCount64(uint64(bb))
}

func (bb Bitboard) LSB() int {
	if bb == EmptyBB {
		return -1
	}
	return bits.TrailingZeros64(uint64(bb))
}
