package board

import (
	"testing"
)

func TestBetween(t *testing.T) {
	tests := []struct {
		a    Square
		b    Square
		want Bitboard
	}{
		{A1, H8, 0x40201008040200},
		{A8, H1, 0x2040810204000},
		{D1, D8, 0x8080808080800},
		{H3, D7, 0x102040000000},
		{A2, G2, 0x3e00},
		{D2, A2, 0x600},
		{D4, E5, EmptyBB},
		{D4, D5, EmptyBB},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Between(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("from %v to %v\ngot:\n%v\nwant:\n%v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestExtending(t *testing.T) {
	tests := []struct {
		a    Square
		b    Square
		want Bitboard
	}{
		{A1, H8, 0x40201008040200},
		{A8, H1, 0x2040810204000},
		{D2, D7, 0x800080808080008},
		{H3, D7, 0x400102040000000},
		{A2, G2, 0xbe00},
		{D2, A2, 0xf600},
		{D4, E5, 0x8040200000040201},
		{D4, D5, 0x808080000080808},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Extending(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("from %v to %v\ngot:\n%v\nwant:\n%v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestBeyond(t *testing.T) {
	tests := []struct {
		a    Square
		b    Square
		want Bitboard
	}{
		{A1, H8, EmptyBB},
		{A8, H1, EmptyBB},
		{D2, D7, D8.Bit()},
		{H3, F5, 0x408100000000000},
		{C2, E2, 0xe000},
		{D2, A2, EmptyBB},
		{D4, E5, 0x8040200000000000},
		{D4, D5, 0x808080000000000},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Beyond(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("from %v to %v\ngot:\n%v\nwant:\n%v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDiagonal(t *testing.T) {
	tests := []struct {
		square    Square
		direction DiagonalDirection
		want      Bitboard
	}{
		{D4, NorthWest, C5.Bit() | B6.Bit() | A7.Bit()},
		{D4, NorthEast, E5.Bit() | F6.Bit() | G7.Bit() | H8.Bit()},
		{D4, SouthEast, E3.Bit() | F2.Bit() | G1.Bit()},
		{D4, SouthWest, C3.Bit() | B2.Bit() | A1.Bit()},
		{A1, NorthWest, EmptyBB},
		{A1, SouthEast, EmptyBB},
		{A1, SouthWest, EmptyBB},
		{A1, NorthEast, 0x8040201008040200},
		{B8, NorthWest, EmptyBB},
		{H8, NorthEast, EmptyBB},
		{H8, SouthEast, EmptyBB},
		{H8, SouthWest, 0x40201008040201},
		{B8, SouthWest, A7.Bit()},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Diagonal(tt.square, tt.direction)
			if got != tt.want {
				t.Fatalf("from %v in direction %v\ngot:\n%v\nwant:\n%v", tt.square, tt.direction, got, tt.want)
			}
		})
	}
}
