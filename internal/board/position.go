package board

import (
	"errors"
	"fmt"
	"strings"
)

type CastlingRights uint8

const (
	BlackKingside CastlingRights = 1 << iota
	BlackQueenside
	WhiteKingside
	WhiteQueenside

	NoCastling  CastlingRights = 0
	AllCastling CastlingRights = BlackKingside | BlackQueenside | WhiteKingside | WhiteQueenside
)

var (
	ErrInvalidCastlingLength = errors.New("castling: string must be 1 to 4 characters, or \"-\" for none")
	ErrInvalidCastlingChar   = errors.New("castling: character must be one of 'K', 'Q', 'k', 'q'")
	ErrDuplicateCastlingChar = errors.New("castling: character appears more than once")
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

func ParseCastlingRights(s string) (CastlingRights, error) {
	n := len(s)
	if n < 1 || n > 4 {
		return NoCastling, fmt.Errorf("%w: %q", ErrInvalidCastlingLength, s)
	}
	if s == "-" {
		return NoCastling, nil
	}

	cr := NoCastling
	for _, char := range s {
		last := cr
		switch char {
		case 'k':
			cr |= BlackKingside
		case 'q':
			cr |= BlackQueenside
		case 'K':
			cr |= WhiteKingside
		case 'Q':
			cr |= WhiteQueenside
		default:
			return NoCastling, fmt.Errorf("%w: %q", ErrInvalidCastlingChar, s)
		}
		if last == cr {
			return NoCastling, fmt.Errorf("%w: %q", ErrDuplicateCastlingChar, s)
		}
	}
	return cr, nil
}

type Bitboards struct {
	pieces [6]Bitboard
	colors [2]Bitboard
}

type Position struct {
	Bitboards
	Board          [64]Piece
	SideToMove     Color
	CastlingRights CastlingRights
	EnPassant      Square
	HalfmoveClock  uint8
	Ply            uint16
	Hash           uint64
}

func (pos *Position) PieceBB(piece Piece) Bitboard {
	return pos.pieces[piece.Type] & pos.colors[piece.Color]
}

func (pos *Position) Occupied() Bitboard {
	return pos.colors[White] | pos.colors[Black]
}

func (pos *Position) PlacePiece(piece Piece, sq Square) {
	pos.pieces[piece.Type].SetBit(sq)
	pos.colors[piece.Color].SetBit(sq)
	pos.Board[sq] = piece
}

func (pos *Position) RemovePiece(sq Square) {
	piece, ok := pos.PieceOn(sq)
	if ok {
		pos.pieces[piece.Type].ClearBit(sq)
		pos.colors[piece.Color].ClearBit(sq)
		pos.Board[sq] = NoPiece
	}
}

func (pos *Position) PieceOn(sq Square) (Piece, bool) {
	piece := pos.Board[sq]
	return piece, piece != NoPiece
}

func (pos Position) String() string {
	var sb strings.Builder
	for rank := Rank8; rank >= Rank1; rank-- {
		for file := FileA; file <= FileH; file++ {
			sq := NewSquare(file, rank)
			sb.WriteString(pos.Board[sq].String())
			sb.WriteByte(' ')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func NewBoard() [64]Piece {
	var board [64]Piece
	for sq := range board {
		board[sq] = NoPiece
	}
	return board
}

func NewPosition() Position {
	pos := Position{
		SideToMove:     White,
		CastlingRights: AllCastling,
		EnPassant:      NoSquare,
		HalfmoveClock:  0,
		Ply:            0,
		Hash:           69420,
	}
	pos.Board = NewBoard()
	return pos
}
