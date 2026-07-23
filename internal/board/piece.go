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
