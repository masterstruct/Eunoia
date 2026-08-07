package board

import "testing"

func TestNewMove(t *testing.T) {
	tests := []struct {
		name string
		from Square
		to   Square
	}{
		{"e2e4", E2, E4},
		{"a1h8", A1, H8},
		{"g1f3", G1, F3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMove(tt.from, tt.to)
			if m.From() != tt.from || m.To() != tt.to {
				t.Errorf("expected %v->%v but got %v->%v", tt.from, tt.to, m.From(), m.To())
			}
			if m.IsCapture() || m.IsPromo() || m.IsCastle() || m.IsEnPassant() || m.IsDoublePush() {
				t.Errorf("expected quiet move but got flags set: %v", m)
			}
		})
	}
}

func TestNewDoublePush(t *testing.T) {
	m := NewDoublePush(E2, E4)
	if m.From() != E2 || m.To() != E4 {
		t.Errorf("expected e2->e4 but got %v->%v", m.From(), m.To())
	}
	if !m.IsDoublePush() {
		t.Errorf("expected IsDoublePush true")
	}
	if m.IsCapture() || m.IsPromo() || m.IsCastle() || m.IsEnPassant() {
		t.Errorf("expected only double-push flag set, got: %v", m)
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
			if m.IsCapture() || m.IsPromo() || m.IsEnPassant() || m.IsDoublePush() {
				t.Errorf("expected only castle flag set, got: %v", m)
			}
		})
	}
}

func TestNewCapture(t *testing.T) {
	m := NewCapture(D4, E5)
	if m.From() != D4 || m.To() != E5 {
		t.Errorf("expected d4->e5 but got %v->%v", m.From(), m.To())
	}
	if !m.IsCapture() {
		t.Errorf("expected IsCapture true")
	}
	if m.IsPromo() || m.IsCastle() || m.IsEnPassant() || m.IsDoublePush() {
		t.Errorf("expected only capture flag set, got: %v", m)
	}
}

func TestNewEnPassant(t *testing.T) {
	m := NewEnPassant(D5, E6)
	if m.From() != D5 || m.To() != E6 {
		t.Errorf("expected d5->e6 but got %v->%v", m.From(), m.To())
	}
	if !m.IsEnPassant() {
		t.Errorf("expected IsEnPassant true")
	}
	if !m.IsCapture() {
		t.Errorf("expected IsCapture true (en passant is a capture)")
	}
	if m.IsPromo() || m.IsCastle() || m.IsDoublePush() {
		t.Errorf("expected only en-passant+capture flags set, got: %v", m)
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
			if m.IsCapture() || m.IsCastle() || m.IsEnPassant() || m.IsDoublePush() {
				t.Errorf("expected only promo flags set, got: %v", m)
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
			if m.Promo() != tt.pt {
				t.Errorf("expected promo type %v but got %v", tt.pt, m.Promo())
			}
		})
	}
}
