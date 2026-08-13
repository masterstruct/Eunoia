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
		want File
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

func TestFileString(t *testing.T) {
	tests := []struct {
		name string
		f    File
		want string
	}{
		{"file a", FileA, "a"},
		{"file b", FileB, "b"},
		{"file c", FileC, "c"},
		{"file d", FileD, "d"},
		{"file e", FileE, "e"},
		{"file f", FileF, "f"},
		{"file g", FileG, "g"},
		{"file h", FileH, "h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.String()
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestRankString(t *testing.T) {
	tests := []struct {
		name string
		r    Rank
		want string
	}{
		{"rank 1", Rank1, "1"},
		{"rank 2", Rank2, "2"},
		{"rank 3", Rank3, "3"},
		{"rank 4", Rank4, "4"},
		{"rank 5", Rank5, "5"},
		{"rank 6", Rank6, "6"},
		{"rank 7", Rank7, "7"},
		{"rank 8", Rank8, "8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.String()
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestRank(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Rank
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
	negativeInvalid := -5
	tests := []struct {
		name string
		sq   Square
		want bool
	}{
		{"A1 valid", A1, true},
		{"H8 valid", H8, true},
		{"E4 valid", E4, true},
		{"NoSquare invalid", NoSquare, false},
		{"negative invalid", Square(negativeInvalid), false},
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

func TestSquareUp(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Square
	}{
		{"A1 to A2", A1, A2},
		{"D4 to D5", D4, D5},
		{"H7 to H8", H7, H8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.Up(); got != tt.want {
				t.Fatalf("Up() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSquareDown(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Square
	}{
		{"A2 to A1", A2, A1},
		{"D5 to D4", D5, D4},
		{"H8 to H7", H8, H7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.Down(); got != tt.want {
				t.Fatalf("Down() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSquareLeft(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Square
	}{
		{"B1 to A1", B1, A1},
		{"E4 to D4", E4, D4},
		{"H8 to G8", H8, G8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.Left(); got != tt.want {
				t.Fatalf("Left() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSquareRight(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Square
	}{
		{"A1 to B1", A1, B1},
		{"D4 to E4", D4, E4},
		{"G8 to H8", G8, H8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sq.Right(); got != tt.want {
				t.Fatalf("Right() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSquare(t *testing.T) {
	tests := []struct {
		name string
		file File
		rank Rank
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
		{NoSquare, "-"},
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
