package movegen

import (
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func TestRookMask(t *testing.T) {
	var tests = []struct {
		sq   board.Square
		want board.Bitboard
	}{
		{
			sq:   board.A1,
			want: 0x101010101017e,
		},
		{
			sq:   board.G2,
			want: 0x40404040403e00,
		},
		{
			sq:   board.D4,
			want: 0x8080876080800,
		},
		{
			sq:   board.H8,
			want: 0x7e80808080808000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.sq.String(), func(t *testing.T) {
			got := RookMask(tt.sq)

			if got != tt.want {
				t.Errorf("expected %v, but got %v", tt.want, got)
			}
		})
	}
}

func TestBishopMask(t *testing.T) {
	var tests = []struct {
		sq   board.Square
		want board.Bitboard
	}{
		{
			sq:   board.A1,
			want: 0x40201008040200,
		},
		{
			sq:   board.G2,
			want: 0x2040810200000,
		},
		{
			sq:   board.D4,
			want: 0x40221400142200,
		},
		{
			sq:   board.H4,
			want: 0x10204000402000,
		},
		{
			sq:   board.H8,
			want: 0x40201008040200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.sq.String(), func(t *testing.T) {
			got := BishopMask(tt.sq)

			if got != tt.want {
				t.Errorf("expected %v, but got %v", tt.want, got)
			}
		})
	}
}

func TestSubsets(t *testing.T) {
	mask := board.Bitboard(0b100011)

	var got []board.Bitboard
	for subset := range Subsets(mask) {
		got = append(got, subset)
	}

	want := []board.Bitboard{
		0b100011,
		0b100010,
		0b100001,
		0b100000,
		0b000011,
		0b000010,
		0b000001,
		0b000000,
	}

	if len(got) != len(want) {
		t.Fatalf("expected %v subsets, got %v", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %b but got %b", want[i], got[i])
		}
	}
}
