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
		{"A1 file is FileA", A1, FileA},
		{"H1 file is FileH", H1, FileH},
		{"E4 file is FileE", E4, FileE},
		{"A8 file is FileA", A8, FileA},
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
		{"A1 rank is Rank1", A1, Rank1},
		{"H1 rank is Rank1", H1, Rank1},
		{"E4 rank is Rank4", E4, Rank4},
		{"A8 rank is Rank8", A8, Rank8},
		{"H8 rank is Rank8", H8, Rank8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.Rank(); got != tt.want {
				t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want bool
	}{
		{"A1 valid", A1, true},
		{"H8 valid", H8, true},
		{"E4 valid", E4, true},
		{"NoSquare invalid", NoSquare, false},
		{"negative invalid", Square(-5), false},
		{"out of range high invalid", Square(64), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sq.IsValid()
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestNewSquare(t *testing.T) {
	tests := []struct {
		name string
		file int
		rank int
		want Square
	}{
		{"a1", FileA, Rank1, A1},
		{"h1", FileH, Rank1, H1},
		{"a8", FileA, Rank8, A8},
		{"h8", FileH, Rank8, H8},
		{"e4", FileE, Rank4, E4},
		{"d5", FileD, Rank5, D5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSquare(tt.file, tt.rank)
			if got != tt.want {
				t.Errorf("NewSquare(%d, %d) = %d, want %d", tt.file, tt.rank, got, tt.want)
			}
		})
	}
}

func TestNewSquareRoundTrip(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		got := NewSquare(sq.File(), sq.Rank())
		if got != sq {
			t.Errorf("round trip failed: square %d -> file %d, rank %d -> %d",
				sq, sq.File(), sq.Rank(), got)
		}
	}
}
