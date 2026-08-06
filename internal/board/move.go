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
	flagQuiet           Move = 0b0000
	flagDoublePush      Move = 0b0001
	flagCastleQueenside Move = 0b0010
	flagCastleKingside  Move = 0b0011

	flagCapture   Move = 0b1000
	flagEnPassant Move = 0b1001

	flagPromoKnight Move = 0b0100
	flagPromoBishop Move = 0b0101
	flagPromoRook   Move = 0b0110
	flagPromoQueen  Move = 0b0111
)

const NullMove = 0 // from A1 to A1 - impossible!

func newMove(from, to Square, flags Move) Move {
	return flags<<flagShift | Move(to)<<toShift | Move(from)
}

func promoFlag(pt PieceType) Move {
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

func NewDoublePush(from, to Square) Move {
	return newMove(from, to, flagDoublePush)
}

func NewCastle(from, to Square, kingside bool) Move {
	if kingside {
		return newMove(from, to, flagCastleKingside)
	}
	return newMove(from, to, flagCastleQueenside)
}

func NewCapture(from, to Square) Move {
	return newMove(from, to, flagCapture)
}

func NewEnPassant(from, to Square) Move {
	return newMove(from, to, flagEnPassant)
}

func (m Move) From() Square {
	return Square(m & squareMask)
}
func (m Move) To() Square {
	return Square((m >> toShift) & squareMask)
}
func (m Move) IsCapture() bool {
	// TODO: make this nicer
	return (m & 0x8000) != 0
}
func (m Move) IsPromo() bool {
	return false
}
func (m Move) Promo() PieceType {
	return NoPieceType
}
func (m Move) IsCastle() bool {
	// TODO: make this nicer
	return (m & 0xE000) == 0x2000
}
func (m Move) IsEnPassant() bool {
	return (m >> flagShift) == flagEnPassant
}
func (m Move) IsDoublePush() bool {
	return (m >> flagShift) == flagDoublePush
}
func (m Move) IsQuiet() bool {
	return false
}
