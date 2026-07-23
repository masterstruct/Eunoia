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
