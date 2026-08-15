package board

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CastlingRights uint8

type CastlingRookSquares struct {
	WhiteKingside  Square
	WhiteQueenside Square
	BlackKingside  Square
	BlackQueenside Square
}

// please don't write at the same time :)
var CastlingRooks CastlingRookSquares

func (pos *Position) NewCastlingRooks() CastlingRookSquares {
	rs := CastlingRookSquares{NoSquare, NoSquare, NoSquare, NoSquare}

	if pos.CastlingRights.Has(WhiteKingside) {
		rs.WhiteKingside = scanRook(pos.PieceBB(WhiteRook), pos.WhiteKing, +1)
	}
	if pos.CastlingRights.Has(WhiteQueenside) {
		rs.WhiteQueenside = scanRook(pos.PieceBB(WhiteRook), pos.WhiteKing, -1)
	}
	if pos.CastlingRights.Has(BlackKingside) {
		rs.BlackKingside = scanRook(pos.PieceBB(BlackRook), pos.BlackKing, +1)
	}
	if pos.CastlingRights.Has(BlackQueenside) {
		rs.BlackQueenside = scanRook(pos.PieceBB(BlackRook), pos.BlackKing, -1)
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

const (
	ResetColor = "\033[0m"
)

var (
	ErrInvalidCastlingLength = errors.New("castling: string must be 1 to 4 characters, or \"-\" for none")
	ErrInvalidCastlingChar   = errors.New("castling: character must be one of 'K', 'Q', 'k', 'q' or rook file letters")
	ErrDuplicateCastlingChar = errors.New("castling: character appears more than once")
	ErrMixedCastlingNotation = errors.New("castling: cannot mix standard (KQkq) and Shredder-FEN (file letter) notation")
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
		return NoCastling, fmt.Errorf("%w: %q", ErrInvalidCastlingLength, s)
	}
	if n == 1 && s[0] == '-' {
		return NoCastling, nil
	}
	hasStandard := strings.ContainsAny(s, "KQkq")
	hasShredder := strings.ContainsAny(s, "ABCDEFGHabcdefgh")
	if hasStandard && hasShredder {
		return NoCastling, fmt.Errorf("%w: %q", ErrMixedCastlingNotation, s)
	}

	var cr CastlingRights

	// standard KQkq form
	switch s[0] {
	case 'K', 'Q', 'k', 'q':
		for _, char := range s {
			var castling CastlingRights
			switch char {
			case 'k':
				castling = BlackKingside
			case 'q':
				castling = BlackQueenside
			case 'K':
				castling = WhiteKingside
			case 'Q':
				castling = WhiteQueenside
			default:
				return NoCastling, fmt.Errorf("%w: %q", ErrInvalidCastlingChar, s)
			}

			if cr&castling != 0 {
				return NoCastling, fmt.Errorf("%w: %q", ErrDuplicateCastlingChar, s)
			}
			cr |= castling
		}
		return cr, nil
	}

	// shredder form
	for _, char := range s {
		// validate and normalize file
		var file File
		switch {
		case 'A' <= char && char <= 'H':
			file = File(char - 'A')
			kingFile := whiteKingSq.File()

			if file == kingFile {
				return NoCastling, fmt.Errorf("%w: %q", ErrInvalidCastlingChar, s)
			}

			last := cr
			if file < kingFile {
				cr |= WhiteQueenside
			} else {
				cr |= WhiteKingside
			}

			if last == cr {
				return NoCastling, fmt.Errorf("%w: %q", ErrDuplicateCastlingChar, s)
			}
		case 'a' <= char && char <= 'h':
			file = File(char - 'a')
			kingFile := blackKingSq.File()

			if file == kingFile {
				return NoCastling, fmt.Errorf("%w: %q", ErrInvalidCastlingChar, s)
			}

			last := cr
			if file < kingFile {
				cr |= BlackQueenside
			} else {
				cr |= BlackKingside
			}

			if last == cr {
				return NoCastling, fmt.Errorf("%w: %q", ErrDuplicateCastlingChar, s)
			}
		default:
			return NoCastling, fmt.Errorf("%w: %q", ErrInvalidCastlingChar, s)
		}
	}
	return cr, nil
}

type Bitboards struct {
	Pieces [6]Bitboard
	Colors [2]Bitboard
}

type Position struct {
	Bitboards
	Board                [64]Piece
	SideToMove           Color
	CastlingRights       CastlingRights
	EnPassant            Square
	HalfmoveClock        uint8
	Ply                  uint16
	Hash                 uint64
	WhiteKing, BlackKing Square
}

func (pos *Position) PieceBB(piece Piece) Bitboard {
	return pos.Pieces[piece.Type] & pos.Colors[piece.Color]
}

func (pos *Position) Occupied() Bitboard {
	return pos.Colors[White] | pos.Colors[Black]
}

func (pos *Position) PlacePiece(piece Piece, sq Square) {
	pos.Pieces[piece.Type].SetBit(sq)
	pos.Colors[piece.Color].SetBit(sq)
	pos.Board[sq] = piece
}

func (pos *Position) RemovePiece(sq Square) {
	piece, ok := pos.PieceOn(sq)
	if ok {
		pos.Pieces[piece.Type].ClearBit(sq)
		pos.Colors[piece.Color].ClearBit(sq)
		pos.Board[sq] = NoPiece
	}
}

func (pos *Position) PieceOn(sq Square) (Piece, bool) {
	piece := pos.Board[sq]
	return piece, piece != NoPiece
}

func (pos *Position) String() string {
	// If using VSCode, go to settings
	// and set this setting:
	// "terminal.integrated.minimumContrastRatio": 1,
	// this will fix the colors of the pretty board

	var sb strings.Builder
	for rank := Rank8; rank >= Rank1; rank-- {
		// ranks
		sb.WriteString(ForegroundRGB(235, 160, 172))
		sb.WriteString(strconv.Itoa(int(rank + 1)))
		sb.WriteByte(' ')

		for file := FileA; file <= FileH; file++ {
			sq := NewSquare(file, rank)
			sqColor := sq.Color()

			sb.WriteString(squareBG(sqColor))

			piece, ok := pos.PieceOn(sq)
			if ok {
				sb.WriteString(pieceFG(piece.Color, sqColor))
				sb.WriteString(piece.PrettyString())
			} else {
				sb.WriteByte(' ')
			}

			sb.WriteByte(' ')
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

func squareBG(sqColor Color) string {
	if sqColor == Black {
		return bgRGB(152, 124, 180)
	}
	return bgRGB(229, 218, 241)
}

func pieceFG(pieceColor Color, sqColor Color) string {
	if pieceColor == Black {
		return ForegroundRGB(49, 50, 68)
	}
	if sqColor == Black {
		// light piece 1
		return ForegroundRGB(211, 183, 146)
	}
	// light piece 2
	return ForegroundRGB(196, 173, 146)
}

func ForegroundRGB(r, g, b int) string {
	// print pretty RGB colors for text foreground
	return fmt.Sprintf("\033[1;38;2;%d;%d;%dm", r, g, b)
}

func bgRGB(r, g, b int) string {
	// print pretty RBG colors for text background
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
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
