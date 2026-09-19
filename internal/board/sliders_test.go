package board

import (
	"testing"
)

func TestRookAttacks(t *testing.T) {
	tests := []struct {
		name     string
		sq       Square
		occupied Bitboard
		want     Bitboard
	}{
		{
			"empty board center d4",
			D4, EmptyBB,
			bbFrom(
				D1, D2, D3, D5, D6, D7, D8,
				A4, B4, C4, E4, F4, G4, H4,
			),
		},
		{
			"empty board corner a1",
			A1, EmptyBB,
			bbFrom(
				A2, A3, A4, A5, A6, A7, A8,
				B1, C1, D1, E1, F1, G1, H1,
			),
		},
		{
			"empty board corner h8",
			H8, EmptyBB,
			bbFrom(
				H1, H2, H3, H4, H5, H6, H7,
				A8, B8, C8, D8, E8, F8, G8,
			),
		},
		{
			"single blocker north",
			D4, bbFrom(D6),
			bbFrom(
				D1, D2, D3, D5, D6,
				A4, B4, C4, E4, F4, G4, H4,
			),
		},
		{
			"blocker on adjacent square",
			D4, bbFrom(D5),
			bbFrom(
				D1, D2, D3, D5,
				A4, B4, C4, E4, F4, G4, H4,
			),
		},
		{
			"blockers in all four directions",
			D4, bbFrom(D6, D2, B4, F4),
			bbFrom(
				D3, D2,
				D5, D6,
				C4, B4,
				E4, F4,
			),
		},
		{
			"fully surrounded, only adjacent squares",
			D4, bbFrom(D5, D3, C4, E4),
			bbFrom(D5, D3, C4, E4),
		},
		{
			"rook on edge, blocker mid-file",
			A1, bbFrom(A4),
			bbFrom(
				A2, A3, A4,
				B1, C1, D1, E1, F1, G1, H1,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.occupied.SetBit(tt.sq)

			got := RookAttacks(tt.sq, tt.occupied)
			if got != tt.want {
				t.Errorf("RookAttacks(%v, %v):\ngot:\n%v\nwant:\n%v", tt.sq, tt.occupied, got, tt.want)
			}
		})
	}
}

func TestBishopAttacks(t *testing.T) {
	tests := []struct {
		name     string
		sq       Square
		occupied Bitboard
		want     Bitboard
	}{
		{
			"empty board center d4",
			D4, EmptyBB,
			bbFrom(
				A1, B2, C3, E5, F6, G7, H8,
				A7, B6, C5, E3, F2, G1,
			),
		},
		{
			"empty board corner a1",
			A1, EmptyBB,
			bbFrom(B2, C3, D4, E5, F6, G7, H8),
		},
		{
			"empty board corner h1",
			H1, EmptyBB,
			bbFrom(G2, F3, E4, D5, C6, B7, A8),
		},
		{
			"single blocker northeast",
			D4, bbFrom(F6),
			bbFrom(
				A1, B2, C3, E5, F6,
				A7, B6, C5, E3, F2, G1,
			),
		},
		{
			"blockers on all four diagonals",
			D4, bbFrom(F6, B2, F2, B6),
			bbFrom(
				C3, B2,
				E5, F6,
				C5, B6,
				E3, F2,
			),
		},
		{
			"bishop on edge a4, half-diagonals only",
			A4, EmptyBB,
			bbFrom(
				B5, C6, D7, E8,
				B3, C2, D1,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.occupied.SetBit(tt.sq)

			got := BishopAttacks(tt.sq, tt.occupied)
			if got != tt.want {
				t.Errorf("BishopAttacks(%v, %v):\ngot:\n%v\nwant:\n%v", tt.sq, tt.occupied, got, tt.want)
			}
		})
	}
}

func TestQueenAttacks(t *testing.T) {
	tests := []struct {
		name     string
		sq       Square
		occupied Bitboard
	}{
		{"empty board center d4", D4, EmptyBB},
		{"empty board corner a1", A1, EmptyBB},
		{"with blockers", D4, bbFrom(D6, F6, B4, C3)},
		{"fully surrounded", D4, bbFrom(
			D5, D3, C4, E4,
			C5, E5, C3, E3,
		)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.occupied.SetBit(tt.sq)

			want := RookAttacks(tt.sq, tt.occupied) | BishopAttacks(tt.sq, tt.occupied)
			got := QueenAttacks(tt.sq, tt.occupied)
			if got != want {
				t.Errorf("QueenAttacks(%v, %v):\ngot:\n%v\nwant (rook|bishop):\n%v", tt.sq, tt.occupied, got, want)
			}
		})
	}
}

func TestMagicAttacksCorrect(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		mask := RookMask(sq)

		for blockers := range Subsets(mask) {
			got := RookAttacks(sq, blockers)
			want := rookAttacksSlow(sq, blockers)

			if got != want {
				t.Fatalf("square %v blockers %v: got %v want %v",
					sq, blockers, got, want)
			}
		}
	}
}
