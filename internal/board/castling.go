package board

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errInvalidCastlingLength = errors.New("castling: string must be 1 to 4 characters, or \"-\" for none")
	errInvalidCastlingChar   = errors.New("castling: character must be one of 'K', 'Q', 'k', 'q' or rook file letters")
	errDuplicateCastlingChar = errors.New("castling: character appears more than once")
	errInvalidCastlingState  = errors.New("castling: string does not match board position")
)

type CastlingRights uint8

const (
	BlackQueenside CastlingRights = iota
	BlackKingside
	WhiteQueenside
	WhiteKingside
)

// black queenside, black kingside, white queenside, white kingside,
// if the square is 0 (NoSquare), the castling right isn't set.
type Castling [4]Square

var noCastling = Castling{NoSquare, NoSquare, NoSquare, NoSquare}

func (c Castling) IsEmpty() bool {
	for _, sq := range c {
		if sq != NoSquare {
			return false
		}
	}
	return true
}

func (c *Castling) Set(rookSq Square, kingside bool) {
	color := Black
	if rookSq.Rank() == Rank1 {
		color = White
	}

	i := color * 2

	if kingside {
		c[i+1] = rookSq
	} else {
		c[i] = rookSq
	}
}

func (c Castling) Has(right CastlingRights) bool {
	return c[right] != NoSquare
}

func (c Castling) HasColor(color Color) bool {
	cr := CastlingRights(color * 2)
	return c.Has(cr) || c.Has(cr+1)
}

func (c Castling) HasSquare(square Square) (bool, CastlingRights) {
	for i, sq := range c {
		if sq == square {
			return true, CastlingRights(i)
		}
	}
	return false, 5
}

func (c *Castling) Remove(right CastlingRights) {
	c[right] = NoSquare
}

func (c *Castling) Clear(color Color) {
	i := color * 2
	c[i] = NoSquare
	c[i+1] = NoSquare
}

// no castling => 0
// ...
// all castling => 15
func (c Castling) ToIndex() uint8 {
	var idx uint8
	for i, sq := range c {
		if sq != NoSquare {
			idx |= 1 << i
		}
	}
	return idx
}

func (c Castling) String(chess960 bool) string {
	if c.IsEmpty() {
		return "-"
	}

	var b strings.Builder

	write := func(ok bool, s string) {
		if ok {
			b.WriteString(s)
		}
	}

	if chess960 {
		write(c.Has(WhiteQueenside), strings.ToUpper(c[WhiteQueenside].File().String()))
		write(c.Has(WhiteKingside), strings.ToUpper(c[WhiteKingside].File().String()))
		write(c.Has(BlackQueenside), c[BlackQueenside].File().String())
		write(c.Has(BlackKingside), c[BlackKingside].File().String())
	} else {
		write(c.Has(WhiteKingside), "K")
		write(c.Has(WhiteQueenside), "Q")
		write(c.Has(BlackKingside), "k")
		write(c.Has(BlackQueenside), "q")
	}

	return b.String()
}

func ParseCastlingRights(s string, whiteKingSq, blackKingSq Square, whiteRooks, blackRooks Bitboard) (Castling, error) {
	n := len(s)

	if n == 1 && s[0] == '-' {
		return noCastling, nil
	}
	if n == 0 || n > 4 {
		return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingLength, s)
	}

	rights := noCastling

	var color Color
	var rooks Bitboard
	var rookSq Square
	var kingSq Square

	for _, char := range s {
		if char&0x20 == 0 {
			// uppercase -> White
			color = White
			kingSq = whiteKingSq
			rooks = whiteRooks
		} else {
			// lowercase -> Black
			color = Black
			kingSq = blackKingSq
			rooks = blackRooks
		}
		char |= 0x20 // fast lowercase conversion

		switch {
		case char == 'k', char == 'q':
			rookSq = scanRook(rooks, kingSq, char == 'k')
		case 'a' <= char && char <= 'h':
			file := File(char - 'a')
			rookSq = NewSquare(file, color.ExpectedKingRank())
		default:
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
		}

		if err := validateRookSquare(rookSq, kingSq, rooks, color); err != nil {
			return noCastling, fmt.Errorf("%w: %q", err, s)
		}

		if ok, _ := rights.HasSquare(rookSq); ok {
			return noCastling, fmt.Errorf("%w: %q", errDuplicateCastlingChar, s)
		}

		rights.Set(rookSq, rookSq > kingSq)
	}
	return rights, nil
}

func scanRook(rookBB Bitboard, kingSq Square, kingside bool) Square {
	rank := kingSq.Rank()
	kingFile := kingSq.File()

	if kingside {
		for file := FileH; file > kingFile; file-- {
			sq := NewSquare(file, rank)
			if rookBB.IsBitSet(sq) {
				return sq
			}
		}
	} else {
		for file := range kingFile {
			sq := NewSquare(file, rank)
			if rookBB.IsBitSet(sq) {
				return sq
			}
		}
	}
	return NoSquare
}

func validateRookSquare(rookSq, kingSq Square, rooks Bitboard, color Color) error {
	if rookSq == NoSquare {
		return errInvalidCastlingState
	}
	if !rooks.IsBitSet(rookSq) {
		return errInvalidCastlingState // rook does not exist on board
	}
	if rookSq.Rank() != kingSq.Rank() {
		return errInvalidCastlingState // king and rook on different ranks
	}
	if rookSq == kingSq {
		return errInvalidCastlingState // rook is inside the king..?
	}
	if kingSq.Rank() != color.ExpectedKingRank() {
		return errInvalidCastlingState // king not on back rank
	}
	return nil
}
