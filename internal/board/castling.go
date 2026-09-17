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

	if n == 0 || n > 4 {
		return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingLength, s)
	}
	if n == 1 && s[0] == '-' {
		return noCastling, nil
	}

	rights := noCastling

	var rookSq Square
	var kingSq Square

	for _, char := range s {
		ok := true

		switch {
		case char == 'k':
			kingSq = blackKingSq
			rookSq, ok = scanRook(blackRooks, kingSq, true)
		case char == 'q':
			kingSq = blackKingSq
			rookSq, ok = scanRook(blackRooks, kingSq, false)
		case char == 'K':
			kingSq = whiteKingSq
			rookSq, ok = scanRook(whiteRooks, kingSq, true)
		case char == 'Q':
			kingSq = whiteKingSq
			rookSq, ok = scanRook(whiteRooks, kingSq, false)
		case 'A' <= char && char <= 'H':
			// white
			rookSq = NewSquare(File(char-'A'), Rank1)
			kingSq = whiteKingSq
		case 'a' <= char && char <= 'h':
			// black
			rookSq = NewSquare(File(char-'a'), Rank8)
			kingSq = blackKingSq
		default:
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
		}

		if !ok {
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingState, s)
		}

		if ok, _ := rights.HasSquare(rookSq); ok {
			return noCastling, fmt.Errorf("%w: %q", errDuplicateCastlingChar, s)
		}

		if !(whiteRooks | blackRooks).IsBitSet(rookSq) {
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingState, s)
		}

		if rookSq == kingSq {
			// rook is inside the king..?
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingState, s)
		}

		if kingSq.Rank() != Rank1 && kingSq.Rank() != Rank8 {
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingState, s)
		}

		rights.Set(rookSq, rookSq > kingSq)
	}
	return rights, nil
}

func scanRook(rookBB Bitboard, kingSq Square, kingside bool) (Square, bool) {
	rank := kingSq.Rank()

	startingFile := FileA
	dir := File(1)
	if kingside {
		startingFile = FileH
		dir = -1
	}

	for file := startingFile; file >= FileA && file <= FileH; file += dir {
		sq := NewSquare(file, rank)
		if sq == kingSq {
			return NoSquare, false
		}
		if rookBB.IsBitSet(sq) {
			return sq, true
		}
	}
	return NoSquare, false
}
