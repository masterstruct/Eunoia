package board

// Thank you, 87flowers, for
// your chess move flag scheme
// https://87flowers.com/chess-moveflags/

type Move uint16

const (
	flagQuiet           uint8 = 0b0000
	flagDoublePush      uint8 = 0b0001
	flagCastleQueenside uint8 = 0b0010
	flagCastleKingside  uint8 = 0b0011

	flagCapture   uint8 = 0b1000
	flagEnPassant uint8 = 0b1001

	flagPromoKnight uint8 = 0b0100
	flagPromoBishop uint8 = 0b0101
	flagPromoRook   uint8 = 0b0110
	flagPromoQueen  uint8 = 0b0111
)

const NullMove = 0 // from A1 to A1 - impossible!

func newMove(from, to Square, flags uint8) Move {
	return Move(uint16(flags)<<12 | uint16(to)<<6 | uint16(from))
}
