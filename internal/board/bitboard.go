package board

import (
	"math/bits"
	"strconv"
	"strings"
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
	// if bb == 0, returns 64 (NoSquare)
	return Square(bits.TrailingZeros64(uint64(bb)))
}

func (bb *Bitboard) PopLSB() Square {
	sq := bb.LSB()
	*bb &= *bb - 1
	return sq
}

func (bb Bitboard) String() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatUint(uint64(bb), 10))

	sb.WriteByte('\n')

	for rank := Rank8; rank >= Rank1; rank-- {
		// ranks
		sb.WriteString(ForegroundRGB(250, 179, 135))
		sb.WriteString(strconv.Itoa(int(rank + 1)))
		sb.WriteByte(' ')

		for file := FileA; file <= FileH; file++ {
			sq := NewSquare(file, rank)
			if bb.IsBitSet(sq) {
				sb.WriteString(ForegroundRGB(166, 227, 161))
				sb.WriteString("1 ")
			} else {
				sb.WriteString(ForegroundRGB(243, 139, 168))
				sb.WriteString(". ")
			}
			sb.WriteString(ResetColor)
		}

		sb.WriteByte('\n')
	}

	// files
	sb.WriteString(ForegroundRGB(116, 199, 236))
	sb.WriteString("  a b c d e f g h\n")
	sb.WriteString(ResetColor)

	return sb.String()
}
