package board

type Color uint8

const (
	Black Color = iota
	White
	NoColor
)

func (c Color) String() string {
	switch c {
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

func (c Color) Opponent() Color {
	switch c {
	case Black:
		return White
	case White:
		return Black
	default:
		return NoColor
	}
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

func (pt PieceType) String() string {
	switch pt {
	case Pawn:
		return "p"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Rook:
		return "r"
	case Queen:
		return "q"
	case King:
		return "k"
	default:
		return "."
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

func PieceTypes() [6]PieceType {
	return [6]PieceType{Pawn, Knight, Bishop, Rook, Queen, King}
}

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

func (p Piece) String() string {
	letter := p.Type.String()
	if p.Color == White && p.Type != NoPieceType {
		// capitalize letter
		letter = string(byte(letter[0]) &^ 0x20)
	}
	return letter
}

func NewPiece(pt PieceType, color Color) Piece {
	if pt == NoPieceType || color == NoColor {
		return Piece{NoPieceType, NoColor}
	}
	return Piece{pt, color}
}
