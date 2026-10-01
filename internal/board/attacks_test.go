package board

import (
	"testing"
)

func bbFrom(squares ...Square) Bitboard {
	var bb Bitboard
	for _, sq := range squares {
		bb.SetBit(sq)
	}
	return bb
}

func TestKnightAttacksFrom(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Bitboard
	}{
		{"corner a1", A1, bbFrom(B3, C2)},
		{"corner h1", H1, bbFrom(F2, G3)},
		{"corner a8", A8, bbFrom(B6, C7)},
		{"corner h8", H8, bbFrom(F7, G6)},
		{"almost corner g7", G7, bbFrom(
			E8, E6, F5, H5,
		)},
		{"edge a4", A4, bbFrom(B6, C5, C3, B2)},
		{"between center and edge c3", C3, bbFrom(
			D5, E4, E2, D1,
			B5, A4, A2, B1,
		)},
		{"between center and edge b3", B3, bbFrom(
			A5, C5, D4,
			D2, C1, A1,
		)},
		{"center d4", D4, bbFrom(
			E6, F5, F3, E2,
			C2, B3, B5, C6,
		)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := knightAttacksFrom(tt.sq)
			if got != tt.want {
				t.Errorf("knightAttacksFrom(%v):\ngot:\n%v\nwant:\n%v", tt.sq, got, tt.want)
			}
		})
	}
}

func TestKingAttacksFrom(t *testing.T) {
	tests := []struct {
		name string
		sq   Square
		want Bitboard
	}{
		{"corner a1", A1, bbFrom(A2, B1, B2)},
		{"corner h8", H8, bbFrom(G8, G7, H7)},
		{"edge a4", A4, bbFrom(A3, A5, B3, B4, B5)},
		{"edge h4", H4, bbFrom(H3, G3, G4, G5, H5)},
		{"center d4", D4, bbFrom(
			C3, C4, C5,
			D3, D5,
			E3, E4, E5,
		)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kingAttacksFrom(tt.sq)
			if got != tt.want {
				t.Errorf("kingAttacksFrom(%v):\ngot:\n%v\nwant:\n%v", tt.sq, got, tt.want)
			}
		})
	}
}

func TestPawnAttacksFrom(t *testing.T) {
	tests := []struct {
		name  string
		sq    Square
		color Color
		want  Bitboard
	}{
		{"white corner a1", A1, White, bbFrom(B2)},
		{"white corner h1", H1, White, bbFrom(G2)},
		{"white edge a4", A4, White, bbFrom(B5)},
		{"white edge h4", H4, White, bbFrom(G5)},
		{"white center d4", D4, White, bbFrom(C5, E5)},

		{"black corner a8", A8, Black, bbFrom(B7)},
		{"black corner h8", H8, Black, bbFrom(G7)},
		{"black edge a5", A5, Black, bbFrom(B4)},
		{"black edge h5", H5, Black, bbFrom(G4)},
		{"black center d5", D5, Black, bbFrom(C4, E4)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pawnAttacksFrom(tt.sq, tt.color)
			if got != tt.want {
				t.Errorf("pawnAttacksFrom(%v, %v):\ngot:\n%v\nwant:\n%v", tt.sq, tt.color, got, tt.want)
			}
		})
	}
}

func TestPawnAttacks(t *testing.T) {
	tests := []struct {
		color Color
		bb    Bitboard
		want  Bitboard
	}{
		{White, EmptyBB, EmptyBB},
		{White, 0xff00, 0xff0000},
		{Black, 0x100000, 0x2800},
		{White, 0x81000e700, 0x142800ff0000},
		{Black, 0x40208004051000, 0xa050400a0a28},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := pawnAttacksSetwise(tt.bb, tt.color)
			if got != tt.want {
				t.Errorf("want %v but got %v", tt.want, got)
			}
		})
	}
}
