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
