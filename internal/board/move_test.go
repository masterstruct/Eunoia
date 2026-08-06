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
			if m.From() != tt.from {
				t.Errorf("expected From %v but got %v", tt.from, m.From())
			}
			if m.To() != tt.to {
				t.Errorf("expected To %v but got %v", tt.to, m.To())
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
		t.Errorf("expected e2e4 but got %v%v", m.From(), m.To())
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
