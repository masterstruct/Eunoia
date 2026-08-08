package movegen

import (
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func TestRookAttacks(t *testing.T) {
	tests := []struct {
		name     string
		sq       board.Square
		occupied board.Bitboard
		want     board.Bitboard
	}{
		{
			"empty board center d4",
			board.D4, board.EmptyBB,
			bbFrom(
				board.D1, board.D2, board.D3, board.D5, board.D6, board.D7, board.D8,
				board.A4, board.B4, board.C4, board.E4, board.F4, board.G4, board.H4,
			),
		},
		{
			"empty board corner a1",
			board.A1, board.EmptyBB,
			bbFrom(
				board.A2, board.A3, board.A4, board.A5, board.A6, board.A7, board.A8,
				board.B1, board.C1, board.D1, board.E1, board.F1, board.G1, board.H1,
			),
		},
		{
			"empty board corner h8",
			board.H8, board.EmptyBB,
			bbFrom(
				board.H1, board.H2, board.H3, board.H4, board.H5, board.H6, board.H7,
				board.A8, board.B8, board.C8, board.D8, board.E8, board.F8, board.G8,
			),
		},
		{
			"single blocker north",
			board.D4, bbFrom(board.D6),
			bbFrom(
				board.D1, board.D2, board.D3, board.D5, board.D6,
				board.A4, board.B4, board.C4, board.E4, board.F4, board.G4, board.H4,
			),
		},
		{
			"blocker on adjacent square",
			board.D4, bbFrom(board.D5),
			bbFrom(
				board.D1, board.D2, board.D3, board.D5,
				board.A4, board.B4, board.C4, board.E4, board.F4, board.G4, board.H4,
			),
		},
		{
			"blockers in all four directions",
			board.D4, bbFrom(board.D6, board.D2, board.B4, board.F4),
			bbFrom(
				board.D3, board.D2,
				board.D5, board.D6,
				board.C4, board.B4,
				board.E4, board.F4,
			),
		},
		{
			"fully surrounded, only adjacent squares",
			board.D4, bbFrom(board.D5, board.D3, board.C4, board.E4),
			bbFrom(board.D5, board.D3, board.C4, board.E4),
		},
		{
			"rook on edge, blocker mid-file",
			board.A1, bbFrom(board.A4),
			bbFrom(
				board.A2, board.A3, board.A4,
				board.B1, board.C1, board.D1, board.E1, board.F1, board.G1, board.H1,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RookAttacks(tt.sq, tt.occupied)
			if got != tt.want {
				t.Errorf("RookAttacks(%v, %v):\ngot:\n%v\nwant:\n%v", tt.sq, tt.occupied, got, tt.want)
			}
		})
	}
}
