package board

import "math/bits"

type Bitboard uint64

const FullBB Bitboard = ^Bitboard(0)
const EmptyBB Bitboard = 0

var RankBB = [8]Bitboard{
	0x00000000000000FF,
	0x000000000000FF00,
	0x0000000000FF0000,
	0x00000000FF000000,
	0x000000FF00000000,
	0x0000FF0000000000,
	0x00FF000000000000,
	0xFF00000000000000,
}

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
	if bit != NoSquare {
		bb.ClearBit(bit)
	}
	return bit
}
