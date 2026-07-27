package board

import (
	"errors"
	"testing"
)

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

func TestColor(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Color
	}{
		{"A1 is dark", A1, Black},
		{"B1 is light", B1, White},
		{"D4 is dark", D4, Black},
		{"E4 is light", E4, White},
		{"H8 is dark", H8, Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.Color(); got != tt.want {
				t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
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

func TestString(t *testing.T) {
	tests := []struct {
		sq   Square
		want string
	}{
		{A1, "a1"},
		{H1, "h1"},
		{A8, "a8"},
		{H8, "h8"},
		{E4, "e4"},
		{D5, "d5"},
		{B2, "b2"},
		{G7, "g7"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.sq.String()
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestParseSquare(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Square
		wantErr error
	}{
		{"a1", "a1", A1, nil},
		{"h1", "h1", H1, nil},
		{"a8", "a8", A8, nil},
		{"h8", "h8", H8, nil},
		{"e4", "e4", E4, nil},
		{"d5", "d5", D5, nil},

		{"empty string", "", NoSquare, ErrInvalidSquareLength},
		{"too short", "e", NoSquare, ErrInvalidSquareLength},
		{"too long", "e44", NoSquare, ErrInvalidSquareLength},

		{"invalid file low", "`4", NoSquare, ErrInvalidFile},
		{"invalid file high", "i4", NoSquare, ErrInvalidFile},
		{"invalid file digit", "1a", NoSquare, ErrInvalidFile},
		{"uppercase file", "E4", E4, nil},

		{"invalid rank zero", "e0", NoSquare, ErrInvalidRank},
		{"invalid rank nine", "e9", NoSquare, ErrInvalidRank},
		{"invalid rank letter", "ee", NoSquare, ErrInvalidRank},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSquare(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v but got %v", tt.wantErr, err)
				}
				if got != NoSquare {
					t.Errorf("expected NoSquare on error but got %v", got)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
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
