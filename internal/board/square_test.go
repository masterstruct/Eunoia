package board

import "testing"

func TestSquareConstants(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Square
	}{
		{"A1 is 0", A1, 0},
		{"H1 is 7", H1, 7},
		{"A8 is 56", A8, 56},
		{"H8 is 63", H8, 63},
		{"E4 is 28", E4, 28},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sq != tt.want {
				t.Errorf("%s: got %d, want %d", tt.name, tt.sq, tt.want)
			}
		})
	}
}

func TestFile(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want int
	}{
		{"A1 file is 0", A1, 0},
		{"H1 file is 7", H1, 7},
		{"E4 file is 4", E4, 4},
		{"A8 file is 0", A8, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.File(); got != tt.want {
				t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

func TestRank(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want int
	}{
		{"A1 rank is 0", A1, 0},
		{"H1 rank is 0", H1, 0},
		{"E4 rank is 3", E4, 3},
		{"A8 rank is 7", A8, 7},
		{"H8 rank is 7", H8, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.Rank(); got != tt.want {
				t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}
