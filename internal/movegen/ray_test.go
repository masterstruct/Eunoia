package movegen

import (
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func TestBetween(t *testing.T) {
	tests := []struct {
		a    board.Square
		b    board.Square
		want board.Bitboard
	}{
		{board.A1, board.H8, 0x40201008040200},
		{board.A8, board.H1, 0x2040810204000},
		{board.D1, board.D8, 0x8080808080800},
		{board.H3, board.D7, 0x102040000000},
		{board.A2, board.G2, 0x3e00},
		{board.D2, board.A2, 0x600},
		{board.D4, board.E5, board.EmptyBB},
		{board.D4, board.D5, board.EmptyBB},
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
		a    board.Square
		b    board.Square
		want board.Bitboard
	}{
		{board.A1, board.H8, 0x40201008040200},
		{board.A8, board.H1, 0x2040810204000},
		{board.D2, board.D7, 0x800080808080008},
		{board.H3, board.D7, 0x400102040000000},
		{board.A2, board.G2, 0xbe00},
		{board.D2, board.A2, 0xf600},
		{board.D4, board.E5, 0x8040200000040201},
		{board.D4, board.D5, 0x808080000080808},
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
