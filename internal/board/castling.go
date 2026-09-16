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
	errMixedCastlingNotation = errors.New("castling: cannot mix standard (KQkq) and Shredder-FEN (file letter) notation")
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

func ParseCastlingRights(s string, whiteKingSq, blackKingSq Square) (Castling, error) {
	// TODO: add X-fen support
	n := len(s)
	if n == 0 || n > 4 {
		return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingLength, s)
	}
	if n == 1 && s[0] == '-' {
		return noCastling, nil
	}
	hasStandard := strings.ContainsAny(s, "KQkq")
	hasShredder := strings.ContainsAny(s, "ABCDEFGHabcdefgh")
	if hasStandard && hasShredder {
		return noCastling, fmt.Errorf("%w: %q", errMixedCastlingNotation, s)
	}

	rights := noCastling
	var sq Square
	var kingside bool

	// standard KQkq form
	switch s[0] {
	case 'K', 'Q', 'k', 'q':
		for _, char := range s {
			switch char {
			case 'k':
				sq = H8
				kingside = true
			case 'q':
				sq = A8
				kingside = false
			case 'K':
				sq = H1
				kingside = true
			case 'Q':
				sq = A1
				kingside = false
			default:
				return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
			}

			if ok, _ := rights.HasSquare(sq); ok {
				return noCastling, fmt.Errorf("%w: %q", errDuplicateCastlingChar, s)
			}
			rights.Set(sq, kingside)
		}
		return rights, nil
	}

	// shredder form
	var file File
	var kingFile File

	for _, char := range s {
		// validate and normalize file

		switch {
		case 'A' <= char && char <= 'H':
			// white
			file = File(char - 'A')
			rookSq := NewSquare(file, Rank1)

			if ok, _ := rights.HasSquare(rookSq); ok {
				return noCastling, fmt.Errorf("%w: %q", errDuplicateCastlingChar, s)
			}

			rights.Set(rookSq, rookSq > whiteKingSq)
			kingFile = whiteKingSq.File()
		case 'a' <= char && char <= 'h':
			// black
			file = File(char - 'a')
			rookSq := NewSquare(file, Rank8)

			if ok, _ := rights.HasSquare(rookSq); ok {
				return noCastling, fmt.Errorf("%w: %q", errDuplicateCastlingChar, s)
			}

			rights.Set(rookSq, rookSq > blackKingSq)
			kingFile = blackKingSq.File()
		default:
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
		}

		if file == kingFile {
			// rook is inside the king..?
			return noCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
		}
	}
	return rights, nil
}

// TODO: X-fens
func FindCastlingRooks(pos *Position) Castling {
	castling := noCastling

	kingSq := pos.KingSq[White]
	if pos.Castling.Has(WhiteQueenside) {
		castling.Set(scanRook(pos.PieceBB(WhiteRook), kingSq, -1), false)
	}
	if pos.Castling.Has(WhiteKingside) {
		castling.Set(scanRook(pos.PieceBB(WhiteRook), kingSq, +1), true)
	}

	kingSq = pos.KingSq[Black]
	if pos.Castling.Has(BlackQueenside) {
		castling.Set(scanRook(pos.PieceBB(BlackRook), kingSq, -1), false)
	}
	if pos.Castling.Has(BlackKingside) {
		castling.Set(scanRook(pos.PieceBB(BlackRook), kingSq, +1), true)
	}

	return castling
}

func scanRook(rookBB Bitboard, kingSq Square, dir File) Square {
	rank := kingSq.Rank()
	for file := kingSq.File() + dir; file >= FileA && file <= FileH; file += dir {
		sq := NewSquare(file, rank)
		if rookBB.IsBitSet(sq) {
			return sq
		}
	}
	return NoSquare
}
