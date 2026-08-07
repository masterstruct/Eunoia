package board

import "testing"

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
				t.Errorf("expected quiet move but got flags set: %v", m)
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
				t.Errorf("expected only double push flag set, got: %v", m)
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
		{"white kingside", E1, G1, true},
		{"white queenside", E1, C1, false},
		{"black kingside", E8, G8, true},
		{"black queenside", E8, C8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewCastle(tt.from, tt.to, tt.kingside)
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
				t.Errorf("expected only castle flag set, got: %v", m)
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
				t.Errorf("expected only capture flag set, got: %v", m)
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
				t.Errorf("expected only en passant flag set, got: %v", m)
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
				t.Errorf("expected only promo flag set, got: %v", m)
			}
		})
	}
}

func TestNewPromoInvalidPieceType(t *testing.T) {
	// edge case: bogus PieceType input falls through promoFlag's default branch
	m := NewPromo(E7, E8, Pawn)
	if m.IsPromo() {
		t.Errorf("expected invalid promo piece type to NOT set promo flag, got: %v", m)
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

func TestNullMove(t *testing.T) {
	var m Move // zero value
	if m != NullMove {
		t.Errorf("expected zero value Move to equal NullMove")
	}
	if m.From() != A1 || m.To() != A1 {
		t.Errorf("expected NullMove to decode as a1->a1, got %v->%v", m.From(), m.To())
	}
}
