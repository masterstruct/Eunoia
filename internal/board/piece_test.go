package board

import "testing"

func TestColorString(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{"white", White, "w"},
		{"black", Black, "b"},
		{"no color", NoColor, "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.String()
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name  string
		input byte
		want  Color
	}{
		{"white lower", 'w', White},
		{"white upper", 'W', White},
		{"black lower", 'b', Black},
		{"black upper", 'B', Black},
		{"invalid", '-', NoColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseColor(tt.input); got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestOpponent(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  Color
	}{
		{"white opponent is black", White, Black},
		{"black opponent is white", Black, White},
		{"no color opponent is no color", NoColor, NoColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.Opponent()
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestPieceTypeString(t *testing.T) {
	tests := []struct {
		name string
		pt   PieceType
		want string
	}{
		{"pawn", Pawn, "p"},
		{"knight", Knight, "n"},
		{"bishop", Bishop, "b"},
		{"rook", Rook, "r"},
		{"queen", Queen, "q"},
		{"king", King, "k"},
		{"no piece type", NoPieceType, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pt.String()
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestParsePieceType(t *testing.T) {
	tests := []struct {
		name  string
		input byte
		want  PieceType
	}{
		{"pawn lower", 'p', Pawn},
		{"pawn upper", 'P', Pawn},
		{"knight lower", 'n', Knight},
		{"knight upper", 'N', Knight},
		{"bishop lower", 'b', Bishop},
		{"bishop upper", 'B', Bishop},
		{"rook lower", 'r', Rook},
		{"rook upper", 'R', Rook},
		{"queen lower", 'q', Queen},
		{"queen upper", 'Q', Queen},
		{"king lower", 'k', King},
		{"king upper", 'K', King},
		{"invalid", '-', NoPieceType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePieceType(tt.input); got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestPieceTypes(t *testing.T) {
	got := PieceTypes()
	want := []PieceType{Pawn, Knight, Bishop, Rook, Queen, King}

	if len(got) != len(want) {
		t.Fatalf("expected length %d but got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("at index %d expected %v but got %v", i, want[i], got[i])
		}
	}
}

func TestPieceString(t *testing.T) {
	tests := []struct {
		name string
		p    Piece
		want string
	}{
		{"white pawn", WhitePawn, "P"},
		{"white knight", WhiteKnight, "N"},
		{"white bishop", WhiteBishop, "B"},
		{"white rook", WhiteRook, "R"},
		{"white queen", WhiteQueen, "Q"},
		{"white king", WhiteKing, "K"},
		{"black pawn", BlackPawn, "p"},
		{"black knight", BlackKnight, "n"},
		{"black bishop", BlackBishop, "b"},
		{"black rook", BlackRook, "r"},
		{"black queen", BlackQueen, "q"},
		{"black king", BlackKing, "k"},
		{"no piece", NoPiece, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.String()
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestNewPiece(t *testing.T) {
	tests := []struct {
		name  string
		pt    PieceType
		color Color
		want  Piece
	}{
		{"white pawn", Pawn, White, WhitePawn},
		{"white king", King, White, WhiteKing},
		{"black pawn", Pawn, Black, BlackPawn},
		{"black king", King, Black, BlackKing},
		{"no piece type", NoPieceType, White, NoPiece},
		{"no color", Pawn, NoColor, NoPiece},
		{"both invalid", NoPieceType, NoColor, NoPiece},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPiece(tt.pt, tt.color)
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}
