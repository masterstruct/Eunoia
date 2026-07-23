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

func ParseColor(s byte) Color {
	switch s {
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
		return ""
	}
}
