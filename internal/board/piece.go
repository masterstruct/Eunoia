package board

type Color uint8

const (
	Black Color = iota
	White
	NoColor
)

func (color Color) String() string {
	switch color {
	case Black:
		return "b"
	case White:
		return "w"
	default:
		return "-"
	}
}

func ParseColor(b byte) Color {
	switch b {
	case 'b', 'B':
		return Black
	case 'w', 'W':
		return White
	default:
		return NoColor
	}
}

func (color Color) Opponent() Color {
	// pls don't do NoColor.Opponent()
	return color ^ 1
}

type PieceType uint8

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
	NoPieceType
)

func (pt PieceType) String() byte {
	switch pt {
	case Pawn:
		return 'p'
	case Knight:
		return 'n'
	case Bishop:
		return 'b'
	case Rook:
		return 'r'
	case Queen:
		return 'q'
	case King:
		return 'k'
	default:
		return '.'
	}
}

func ParsePieceType(b byte) PieceType {
	switch b {
	case 'p', 'P':
		return Pawn
	case 'n', 'N':
		return Knight
	case 'b', 'B':
		return Bishop
	case 'r', 'R':
		return Rook
	case 'q', 'Q':
		return Queen
	case 'k', 'K':
		return King
	default:
		return NoPieceType
	}
}

var PieceTypes = [6]PieceType{Pawn, Knight, Bishop, Rook, Queen, King}

type Piece struct {
	Type  PieceType
	Color Color
}

var (
	BlackPawn   = Piece{Pawn, Black}
	BlackKnight = Piece{Knight, Black}
	BlackBishop = Piece{Bishop, Black}
	BlackRook   = Piece{Rook, Black}
	BlackQueen  = Piece{Queen, Black}
	BlackKing   = Piece{King, Black}

	WhitePawn   = Piece{Pawn, White}
	WhiteKnight = Piece{Knight, White}
	WhiteBishop = Piece{Bishop, White}
	WhiteRook   = Piece{Rook, White}
	WhiteQueen  = Piece{Queen, White}
	WhiteKing   = Piece{King, White}

	NoPiece = Piece{NoPieceType, NoColor}
)

func (p Piece) String() byte {
	letter := p.Type.String()
	if p.Color == White && p.Type != NoPieceType {
		// capitalize letter
		letter &^= 0x20
	}
	return letter
}

func (p Piece) PrettyString() string {
	var symbol rune
	switch p.Type {
	case Pawn:
		symbol = '♟'
	case Knight:
		symbol = '♞'
	case Bishop:
		symbol = '♝'
	case Rook:
		symbol = '♜'
	case Queen:
		symbol = '♛'
	case King:
		symbol = '♚'
	default:
		return " "
	}
	return string(symbol)
}

func NewPiece(pt PieceType, color Color) Piece {
	if pt == NoPieceType || color == NoColor {
		return Piece{NoPieceType, NoColor}
	}
	return Piece{pt, color}
}
