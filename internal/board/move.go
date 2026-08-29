package board

import (
	"sync/atomic"
)

// TODO: move to an approprate file
var chess960 atomic.Bool

func SetChess960(v bool) {
	chess960.Store(v)
}

func IsChess960() bool {
	return chess960.Load()
}

// Thank you, 87flowers, for
// your chess move flag scheme
// https://87flowers.com/chess-moveflags/

// [ flags ][  to  ][ from ]
// [ 15-12 ][ 11-6 ][ 5-0  ]
type Move uint16

const NullMove Move = 0 // from A1 to A1 - impossible!

const (
	squareMask Move = 0x3F   // bits 0-5
	flagMask   Move = 0xF000 // bits 12-15

	squareBits = 6
	toShift    = squareBits     // 6
	flagShift  = squareBits * 2 // 12
)

const (
	flagQuiet      Move = 0x0000
	flagDoublePush Move = 0x1000

	castleMask          Move = 0xE000
	flagCastlePattern   Move = 0x2000
	flagCastleQueenside Move = 0x2000
	flagCastleKingside  Move = 0x3000

	flagCapture   Move = 0x8000
	flagEnPassant Move = 0x9000

	flagPromo       Move = 0x4000
	flagPromoKnight Move = 0x4000
	flagPromoBishop Move = 0x5000
	flagPromoRook   Move = 0x6000
	flagPromoQueen  Move = 0x7000
)

func newMove(from, to Square, flags Move) Move {
	return flags | Move(to)<<toShift | Move(from)
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

func NewCastle(from, to Square) Move {
	if to > from { // kingside
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

func NewPromo(from, to Square, pt PieceType) Move {
	return newMove(from, to, promoFlag(pt))
}

func NewCapturePromo(from, to Square, pt PieceType) Move {
	return newMove(from, to, promoFlag(pt)|flagCapture)
}

func (m Move) From() Square {
	return Square(m & squareMask)
}

func (m Move) To() Square {
	return Square((m >> toShift) & squareMask)
}

func (m Move) IsCapture() bool {
	return (m & flagCapture) != 0
}

func (m Move) IsPromo() bool {
	return (m & flagPromo) != 0
}

func (m Move) Promo() PieceType {
	switch (m & flagMask) &^ flagCapture {
	case flagPromoKnight:
		return Knight
	case flagPromoBishop:
		return Bishop
	case flagPromoRook:
		return Rook
	case flagPromoQueen:
		return Queen
	default:
		return NoPieceType
	}
}

func (m Move) IsCastle() bool {
	// grab last 3 bits and compare them to 001
	return (m & castleMask) == flagCastlePattern
}

func (m Move) IsKingsideCastle() bool {
	return (m & flagMask) == flagCastleKingside
}

func (m Move) IsEnPassant() bool {
	return (m & flagMask) == flagEnPassant
}

func (m Move) IsDoublePush() bool {
	return (m & flagMask) == flagDoublePush
}

func (m Move) IsQuiet() bool {
	return (m & flagMask) == flagQuiet
}

func (m Move) String() string {
	from, to := m.From(), m.To()
	if m.IsCastle() && !IsChess960() {
		to = FischerRandomToStandardCastling(m)
	}
	return from.String() + to.String() + promoSuffix(m)
}

func promoSuffix(m Move) string {
	if !m.IsPromo() {
		return ""
	}
	return string(m.Promo().String())
}

func FischerRandomToStandardCastling(move Move) Square {
	switch move.From().Rank() {
	case Rank8:
		if move.IsKingsideCastle() {
			return G8
		} else {
			return C8
		}
	case Rank1:
		if move.IsKingsideCastle() {
			return G1
		} else {
			return C1
		}
	default:
		return NoSquare
	}
}

func (m Move) Raw() uint16 {
	return uint16(m)
}

func (m Move) RawFlags() uint16 {
	return m.Raw() >> flagShift
}
