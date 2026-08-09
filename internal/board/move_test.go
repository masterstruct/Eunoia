package board

import "testing"

// TODO: move to an approprate file
func withChess960(t *testing.T) {
	t.Helper()
	SetChess960(true)
	t.Cleanup(func() { SetChess960(false) })
}

func TestNewMove(t *testing.T) {
	tests := []struct {
		name string
		from Square
		to   Square
	}{
		{"e2e3", E2, E3},
		{"a1h8", A1, H8},
		{"g1f3", G1, F3},
		{"d7d6", D7, D6},
		{"g8f6", G8, F6},
		{"c8g4", C8, G4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMove(tt.from, tt.to)
			if m.From() != tt.from || m.To() != tt.to {
				t.Errorf("expected %v->%v but got %v->%v", tt.from, tt.to, m.From(), m.To())
			}
			if !m.IsQuiet() {
				t.Errorf("expected IsQuiet true")
			}
			if m.IsCapture() || m.IsPromo() || m.IsCastle() || m.IsEnPassant() || m.IsDoublePush() {
				t.Errorf("expected quiet move but got flags set: %04b", m.RawFlags())
			}
		})
	}
}

func TestNewDoublePush(t *testing.T) {
	tests := []struct {
		name string
		from Square
		to   Square
	}{
		{"e2e4", E2, E4},
		{"a2a4", A2, A4},
		{"h2h4", H2, H4},
		{"e7e5", E7, E5},
		{"a7a5", A7, A5},
		{"h7h5", H7, H5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewDoublePush(tt.from, tt.to)
			if m.From() != tt.from || m.To() != tt.to {
				t.Errorf("expected %v->%v but got %v->%v", tt.from, tt.to, m.From(), m.To())
			}
			if !m.IsDoublePush() {
				t.Errorf("expected IsDoublePush true")
			}
			if m.IsQuiet() {
				t.Errorf("expected IsQuiet false")
			}
			if m.IsCapture() || m.IsPromo() || m.IsCastle() || m.IsEnPassant() {
				t.Errorf("expected only double push flag set, got: %04b", m.RawFlags())
			}
		})
	}
}

func TestNewCastle(t *testing.T) {
	tests := []struct {
		name     string
		from, to Square
		kingside bool
	}{
		{"white kingside", E1, H1, true},
		{"white queenside", E1, A1, false},
		{"black kingside", E8, H8, true},
		{"black queenside", E8, A8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewCastle(tt.from, tt.to)
			gotKingside := m.IsKingsideCastle()
			if m.From() != tt.from || m.To() != tt.to {
				t.Errorf("expected %v->%v but got %v->%v", tt.from, tt.to, m.From(), m.To())
			}
			if !m.IsCastle() {
				t.Errorf("expected IsCastle true")
			}
			if m.IsQuiet() {
				t.Errorf("expected IsQuiet false")
			}
			if m.IsCapture() || m.IsPromo() || m.IsEnPassant() || m.IsDoublePush() {
				t.Errorf("expected only castle flag set, got: %04b", m.RawFlags())
			}
			if gotKingside != tt.kingside {
				t.Errorf("expected IsKingsideCastle %v, got %v", tt.kingside, gotKingside)
			}
		})
	}
}

func TestNewCapture(t *testing.T) {
	tests := []struct {
		name string
		from Square
		to   Square
	}{
		{"d4e5", D4, E5},
		{"b4c3", B4, C3},
		{"e2a6", E2, A6},
		{"f6e4", F6, E4},
		{"e5g6", E5, G6},
		{"f3f6", F3, F6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewCapture(tt.from, tt.to)
			if m.From() != tt.from || m.To() != tt.to {
				t.Errorf("expected %v->%v but got %v->%v", tt.from, tt.to, m.From(), m.To())
			}
			if !m.IsCapture() {
				t.Errorf("expected IsCapture true")
			}
			if m.IsQuiet() {
				t.Errorf("expected IsQuiet false")
			}
			if m.IsPromo() || m.IsCastle() || m.IsEnPassant() || m.IsDoublePush() {
				t.Errorf("expected only capture flag set, got: %04b", m.RawFlags())
			}
		})
	}
}

func TestNewEnPassant(t *testing.T) {
	tests := []struct {
		name string
		from Square
		to   Square
	}{
		{"d5c6", D5, C6},
		{"e5d6", E5, D6},
		{"a5b6", A5, B6},
		{"d4c3", D4, C3},
		{"e4d3", E4, D3},
		{"a4b3", A4, B3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewEnPassant(tt.from, tt.to)
			if m.From() != tt.from || m.To() != tt.to {
				t.Errorf("expected %v->%v but got %v->%v", tt.from, tt.to, m.From(), m.To())
			}
			if !m.IsEnPassant() {
				t.Errorf("expected IsEnPassant true")
			}
			if !m.IsCapture() {
				t.Errorf("expected IsCapture true (en passant is a capture)")
			}
			if m.IsQuiet() {
				t.Errorf("expected IsQuiet false")
			}
			if m.IsPromo() || m.IsCastle() || m.IsDoublePush() {
				t.Errorf("expected only en passant flag set, got: %04b", m.RawFlags())
			}
		})
	}
}

func TestNewPromo(t *testing.T) {
	tests := []struct {
		name string
		pt   PieceType
	}{
		{"knight", Knight},
		{"bishop", Bishop},
		{"rook", Rook},
		{"queen", Queen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewPromo(E7, E8, tt.pt)
			if m.From() != E7 || m.To() != E8 {
				t.Errorf("expected e7->e8 but got %v->%v", m.From(), m.To())
			}
			if !m.IsPromo() {
				t.Errorf("expected IsPromo true")
			}
			if m.Promo() != tt.pt {
				t.Errorf("expected promo type %v but got %v", tt.pt, m.Promo())
			}
			if m.IsQuiet() {
				t.Errorf("expected IsQuiet false")
			}
			if m.IsCapture() || m.IsCastle() || m.IsEnPassant() || m.IsDoublePush() {
				t.Errorf("expected only promo flag set, got: %04b", m.RawFlags())
			}
		})
	}
}

func TestNewPromoInvalidPieceType(t *testing.T) {
	// edge case: bogus PieceType input falls through promoFlag's default branch
	m := NewPromo(E7, E8, Pawn)
	if m.IsPromo() {
		t.Errorf("expected invalid promo piece type to NOT set promo flag, got: %04b", m.RawFlags())
	}
}

func TestNewCapturePromo(t *testing.T) {
	tests := []struct {
		name string
		pt   PieceType
	}{
		{"knight", Knight},
		{"bishop", Bishop},
		{"rook", Rook},
		{"queen", Queen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewCapturePromo(D7, C8, tt.pt)
			if m.From() != D7 || m.To() != C8 {
				t.Errorf("expected d7->c8 but got %v->%v", m.From(), m.To())
			}
			if !m.IsCapture() {
				t.Errorf("expected IsCapture true")
			}
			if !m.IsPromo() {
				t.Errorf("expected IsPromo true")
			}
			if m.IsQuiet() {
				t.Errorf("expected IsQuiet false")
			}
			if m.Promo() != tt.pt {
				t.Errorf("expected promo type %v but got %v", tt.pt, m.Promo())
			}
		})
	}
}

func TestNewCapturePromoInvalidPieceType(t *testing.T) {
	m := NewCapturePromo(E7, E8, Pawn)
	if m.IsPromo() {
		t.Errorf("expected invalid promo piece type to NOT set promo flag, got: %04b", m.RawFlags())
	}
	if !m.IsCapture() {
		t.Errorf("expected IsCapture to remain true even with invalid promo type")
	}
}

func TestNullMove(t *testing.T) {
	var m Move // zero value
	if !m.IsNullMove() {
		t.Errorf("expected zero value Move to equal NullMove")
	}
	if m.From() != A1 || m.To() != A1 {
		t.Errorf("expected NullMove to decode as a1->a1, got %v->%v", m.From(), m.To())
	}
}

func TestIsNullMove(t *testing.T) {
	tests := []struct {
		name string
		m    Move
		want bool
	}{
		{"null move", NullMove, true},
		{"quiet move", NewMove(E2, E3), false},
		{"double push", NewDoublePush(E2, E4), false},
		{"capture", NewCapture(E4, D5), false},
		{"en passant", NewEnPassant(E5, D6), false},
		{"kingside castle", NewCastle(E1, H1), false},
		{"queenside castle", NewCastle(E1, A1), false},
		{"promotion", NewPromo(E7, E8, Queen), false},
		{"capture promotion", NewCapturePromo(E7, D8, Queen), false},
		{"same square constructed manually", NewMove(A1, A1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.IsNullMove()
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestMove_String(t *testing.T) {
	t.Run("quiet move", func(t *testing.T) {
		tests := []struct {
			from, to Square
			want     string
		}{
			{E2, E3, "e2e3"},
			{G1, F3, "g1f3"},
			{D7, D6, "d7d6"},
			{C8, G4, "c8g4"},
		}

		for _, tt := range tests {
			m := NewMove(tt.from, tt.to)
			got := m.String()
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		}
	})

	t.Run("double push", func(t *testing.T) {
		tests := []struct {
			from, to Square
			want     string
		}{
			{E2, E4, "e2e4"},
			{A2, A4, "a2a4"},
			{H2, H4, "h2h4"},
			{D7, D5, "d7d5"},
			{A7, A5, "a7a5"},
			{H7, H5, "h7h5"},
		}

		for _, tt := range tests {
			m := NewDoublePush(tt.from, tt.to)
			got := m.String()
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		}
	})

	t.Run("castling", func(t *testing.T) { // NOT chess960
		tests := []struct {
			from, to Square
			want     string
		}{
			{E1, H1, "e1g1"},
			{E1, A1, "e1c1"},
			{E8, H8, "e8g8"},
			{E8, A8, "e8c8"},
		}

		for _, tt := range tests {
			m := NewCastle(tt.from, tt.to)
			got := m.String()
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		}
	})

	t.Run("en passant", func(t *testing.T) {
		tests := []struct {
			from, to Square
			want     string
		}{
			{D5, C6, "d5c6"},
			{E5, D6, "e5d6"},
			{A5, B6, "a5b6"},
			{D4, C3, "d4c3"},
			{E4, D3, "e4d3"},
			{A4, B3, "a4b3"},
		}

		for _, tt := range tests {
			m := NewEnPassant(tt.from, tt.to)
			got := m.String()
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		}
	})

	t.Run("regular promotion", func(t *testing.T) {
		tests := []struct {
			from, to Square
			pt       PieceType
			want     string
		}{
			{D7, D8, Knight, "d7d8n"},
			{C7, C8, Bishop, "c7c8b"},
			{B7, B8, Rook, "b7b8r"},
			{A7, A8, Queen, "a7a8q"},

			{D2, D1, Knight, "d2d1n"},
			{C2, C1, Bishop, "c2c1b"},
			{B2, B1, Rook, "b2b1r"},
			{H2, H1, Queen, "h2h1q"},
		}

		for _, tt := range tests {
			m := NewPromo(tt.from, tt.to, tt.pt)
			got := m.String()
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		}
	})

	t.Run("capture promotion", func(t *testing.T) {
		tests := []struct {
			from, to Square
			pt       PieceType
			want     string
		}{
			{H7, H8, Knight, "h7h8n"},
			{C7, D8, Bishop, "c7d8b"},
			{B7, A8, Rook, "b7a8r"},
			{A7, B8, Queen, "a7b8q"},

			{D2, E1, Knight, "d2e1n"},
			{C2, D1, Bishop, "c2d1b"},
			{B2, A1, Rook, "b2a1r"},
			{H2, G1, Queen, "h2g1q"},
		}

		for _, tt := range tests {
			m := NewCapturePromo(tt.from, tt.to, tt.pt)
			got := m.String()
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		}
	})
}

func TestMoveString_Chess960Castle(t *testing.T) {
	withChess960(t)

	tests := []struct {
		name string
		move Move
		want string
	}{
		{"white kingside", NewCastle(E1, H1), "e1h1"},
		{"white queenside", NewCastle(E1, A1), "e1a1"},
		{"white kingside, rook on f", NewCastle(E1, F1), "e1f1"},
		{"white kingside, rook on c", NewCastle(E1, C1), "e1c1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.move.String(); got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}
