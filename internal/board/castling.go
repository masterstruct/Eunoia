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

type CastlingRookSquares struct {
	WhiteKingside  Square
	WhiteQueenside Square
	BlackKingside  Square
	BlackQueenside Square
}

func (pos *Position) NewCastlingRooks() CastlingRookSquares {
	rs := CastlingRookSquares{NoSquare, NoSquare, NoSquare, NoSquare}

	if pos.CastlingRights.Has(WhiteKingside) {
		rs.WhiteKingside = scanRook(pos.PieceBB(WhiteRook), pos.KingSq[White], +1)
	}
	if pos.CastlingRights.Has(WhiteQueenside) {
		rs.WhiteQueenside = scanRook(pos.PieceBB(WhiteRook), pos.KingSq[White], -1)
	}
	if pos.CastlingRights.Has(BlackKingside) {
		rs.BlackKingside = scanRook(pos.PieceBB(BlackRook), pos.KingSq[Black], +1)
	}
	if pos.CastlingRights.Has(BlackQueenside) {
		rs.BlackQueenside = scanRook(pos.PieceBB(BlackRook), pos.KingSq[Black], -1)
	}
	return rs
}

func scanRook(rookBB Bitboard, kingSq Square, dir File) Square {
	// TODO: use PopLSB
	rank := kingSq.Rank()
	for file := kingSq.File() + dir; file >= FileA && file <= FileH; file += dir {
		sq := NewSquare(file, rank)
		if rookBB.IsBitSet(sq) {
			return sq
		}
	}
	return NoSquare
}

const (
	BlackKingside CastlingRights = 1 << iota
	BlackQueenside
	WhiteKingside
	WhiteQueenside

	NoCastling  CastlingRights = 0
	AllCastling CastlingRights = BlackKingside | BlackQueenside | WhiteKingside | WhiteQueenside
)

func (cr CastlingRights) Has(right CastlingRights) bool {
	return cr&right != NoCastling
}

func (cr *CastlingRights) Remove(right CastlingRights) {
	*cr &^= right
}

func (cr CastlingRights) String() string {
	if cr == NoCastling {
		return "-"
	}

	s := ""
	if cr.Has(WhiteKingside) {
		s += "K"
	}
	if cr.Has(WhiteQueenside) {
		s += "Q"
	}
	if cr.Has(BlackKingside) {
		s += "k"
	}
	if cr.Has(BlackQueenside) {
		s += "q"
	}
	return s
}

func ParseCastlingRights(s string, whiteKingSq, blackKingSq Square) (CastlingRights, error) {
	n := len(s)
	if n == 0 || n > 4 {
		return NoCastling, fmt.Errorf("%w: %q", errInvalidCastlingLength, s)
	}
	if n == 1 && s[0] == '-' {
		return NoCastling, nil
	}
	hasStandard := strings.ContainsAny(s, "KQkq")
	hasShredder := strings.ContainsAny(s, "ABCDEFGHabcdefgh")
	if hasStandard && hasShredder {
		return NoCastling, fmt.Errorf("%w: %q", errMixedCastlingNotation, s)
	}

	var rights CastlingRights

	// standard KQkq form
	switch s[0] {
	case 'K', 'Q', 'k', 'q':
		for _, char := range s {
			var newRights CastlingRights
			switch char {
			case 'k':
				newRights = BlackKingside
			case 'q':
				newRights = BlackQueenside
			case 'K':
				newRights = WhiteKingside
			case 'Q':
				newRights = WhiteQueenside
			default:
				return NoCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
			}

			// newRights already in rights - duplicate entry
			if rights.Has(newRights) {
				return NoCastling, fmt.Errorf("%w: %q", errDuplicateCastlingChar, s)
			}
			rights |= newRights
		}
		return rights, nil
	}

	// shredder form
	var file File
	var kingFile File
	var queenside CastlingRights
	var kingside CastlingRights

	for _, char := range s {
		// validate and normalize file

		switch {
		case 'A' <= char && char <= 'H':
			// white
			file = File(char - 'A')
			queenside = WhiteQueenside
			kingside = WhiteKingside
			kingFile = whiteKingSq.File()
		case 'a' <= char && char <= 'h':
			// black
			file = File(char - 'a')
			queenside = BlackQueenside
			kingside = BlackKingside
			kingFile = blackKingSq.File()
		default:
			return NoCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
		}

		if file == kingFile {
			// rook is inside the king..?
			return NoCastling, fmt.Errorf("%w: %q", errInvalidCastlingChar, s)
		}

		last := rights
		if file < kingFile {
			rights |= queenside
		} else {
			rights |= kingside
		}

		if last == rights {
			return NoCastling, fmt.Errorf("%w: %q", errDuplicateCastlingChar, s)
		}
	}
	return rights, nil
}
