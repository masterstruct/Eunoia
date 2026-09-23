package board

import (
	"math/bits"
)

type Bitboard uint64

const (
	FullBB  Bitboard = ^Bitboard(0)
	EmptyBB Bitboard = 0

	EdgesBB Bitboard = 0xff818181818181ff
)

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

var FileBB = [8]Bitboard{
	0x101010101010101,
	0x202020202020202,
	0x404040404040404,
	0x808080808080808,
	0x1010101010101010,
	0x2020202020202020,
	0x4040404040404040,
	0x8080808080808080,
}

var SquareBB [64]Bitboard

func init() {
	InitBitboards()
}

func InitBitboards() {
	for sq := range 64 {
		SquareBB[sq] = Bitboard(1) << sq
	}
}

func (bb *Bitboard) SetBit(sq Square) {
	*bb |= sq.Bit()
}

func (bb *Bitboard) ClearBit(sq Square) {
	*bb &^= sq.Bit()
}

func (bb Bitboard) IsBitSet(sq Square) bool {
	return (bb & sq.Bit()) != 0
}

func (bb Bitboard) CountBits() int {
	return bits.OnesCount64(uint64(bb))
}

func (bb Bitboard) LSB() Square {
	// if bb == 0, returns 64 (NoSquare)
	return Square(bits.TrailingZeros64(uint64(bb)))
}

func (bb *Bitboard) PopLSB() Square {
	sq := bb.LSB()
	*bb &= *bb - 1
	return sq
}

func (bb Bitboard) northEast() Bitboard {
	return Bitboard(bb << 9 &^ FileA.Bits())
}

func (bb Bitboard) northWest() Bitboard {
	return Bitboard(bb << 7 &^ FileH.Bits())
}

func (bb Bitboard) southEast() Bitboard {
	return Bitboard(bb >> 7 &^ FileA.Bits())
}

func (bb Bitboard) southWest() Bitboard {
	return Bitboard(bb >> 9 &^ FileH.Bits())
}
