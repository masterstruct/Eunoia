package board

type Color uint8

const (
	Black Color = iota
	White
	NoColor
)

func (c Color) ExpectedKingRank() Rank {
	if c == White {
		return Rank1
	}
	return Rank8
}

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
