package board

type PieceType uint8

var PieceTypes = [6]PieceType{Pawn, Knight, Bishop, Rook, Queen, King}

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
	NoPieceType
)

const (
	pawnValue   = 100
	knightValue = 300
	bishopValue = 330
	rookValue   = 500
	queenValue  = 900
)

func (pt PieceType) Value() int32 {
	switch pt {
	case Pawn:
		return pawnValue
	case Knight:
		return knightValue
	case Bishop:
		return bishopValue
	case Rook:
		return rookValue
	case Queen:
		return queenValue
	default:
		return 0
	}
}

type Piece uint8

var Pieces = [12]Piece{BlackPawn, WhitePawn, BlackKnight, WhiteKnight, BlackBishop,
	WhiteBishop, BlackRook, WhiteRook, BlackQueen, WhiteQueen, BlackKing, BlackQueen}

const (
	BlackPawn Piece = iota
	WhitePawn
	BlackKnight
	WhiteKnight
	BlackBishop
	WhiteBishop
	BlackRook
	WhiteRook
	BlackQueen
	WhiteQueen
	BlackKing
	WhiteKing

	NoPiece
)

func NewPiece(pt PieceType, color Color) Piece {
	return Piece(uint8(pt)<<1 | uint8(color))
}

func (p Piece) Type() PieceType {
	return PieceType(p >> 1)
}

func (p Piece) Color() Color {
	return Color(p & 1)
}

func (p Piece) Value() int32 {
	return p.Type().Value()
}

func (p Piece) String() byte {
	letter := p.Type().String()
	if p.Color() == White && p.Type() != NoPieceType {
		// capitalize letter
		letter &^= 0x20
	}
	return letter
}

func (p Piece) PrettyString() string {
	var symbol rune
	switch p.Type() {
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
