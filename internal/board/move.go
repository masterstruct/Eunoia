package board

// Thank you, 87flowers, for
// your chess move flag scheme
// https://87flowers.com/chess-moveflags/

// [ flags ][  to  ][ from ]
// [ 15-12 ][ 11-6 ][ 5-0  ]
type Move uint16

const (
	squareBits      = 6
	squareMask Move = (1 << squareBits) - 1 // 0b111111
	flagMask   Move = 0b1111

	fromShift = 0
	toShift   = squareBits     // 6
	flagShift = squareBits * 2 // 12
)

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

func promoFlag(pt PieceType) uint8 {
	switch pt {
	case Knight:
		return flagPromoKnight
	case Bishop:
		return flagPromoBishop
	case Rook:
		return flagPromoRook
	case Queen:
		return flagPromoQueen
	default:
		return flagQuiet
	}
}

func NewMove(from, to Square) Move {
	return newMove(from, to, flagQuiet)
}

func (m Move) From() Square {
	return Square(m & Move(squareMask))
}
func (m Move) To() Square {
	return NoSquare
}
func (m Move) IsCapture() bool {
	return false
}
func (m Move) IsPromo() bool {
	return false
}
func (m Move) Promo() PieceType {
	return NoPieceType
}
func (m Move) IsCastle() bool {
	return false
}
func (m Move) IsEnPassant() bool {
	return false
}
func (m Move) IsDoublePush() bool {
	return false
}
func (m Move) IsQuiet() bool {
	return false
}
