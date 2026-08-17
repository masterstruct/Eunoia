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
			got := rookMask(tt.sq)

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
			got := bishopMask(tt.sq)

			if got != tt.want {
				t.Errorf("expected %v, but got %v", tt.want, got)
			}
		})
	}
}
