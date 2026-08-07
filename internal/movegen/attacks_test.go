package movegen

import (
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func bbFrom(squares ...board.Square) board.Bitboard {
	var bb board.Bitboard
	for _, sq := range squares {
		bb.SetBit(sq)
	}
	return bb
}

func TestKnightAttacksFrom(t *testing.T) {
	tests := []struct {
		name string
		sq   board.Square
		want board.Bitboard
	}{
		{"corner a1", board.A1, bbFrom(board.B3, board.C2)},
		{"corner h1", board.H1, bbFrom(board.F2, board.G3)},
		{"corner a8", board.A8, bbFrom(board.B6, board.C7)},
		{"corner h8", board.H8, bbFrom(board.F7, board.G6)},
		{"almost corner g7", board.G7, bbFrom(
			board.E8, board.E6, board.F5, board.H5,
		)},
		{"edge a4", board.A4, bbFrom(board.B6, board.C5, board.C3, board.B2)},
		{"between center and edge c3", board.C3, bbFrom(
			board.D5, board.E4, board.E2, board.D1,
			board.B5, board.A4, board.A2, board.B1,
		)},
		{"between center and edge b3", board.B3, bbFrom(
			board.A5, board.C5, board.D4,
			board.D2, board.C1, board.A1,
		)},
		{"center d4", board.D4, bbFrom(
			board.E6, board.F5, board.F3, board.E2,
			board.C2, board.B3, board.B5, board.C6,
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
		sq   board.Square
		want board.Bitboard
	}{
		{"corner a1", board.A1, bbFrom(board.A2, board.B1, board.B2)},
		{"corner h8", board.H8, bbFrom(board.G8, board.G7, board.H7)},
		{"edge a4", board.A4, bbFrom(board.A3, board.A5, board.B3, board.B4, board.B5)},
		{"edge h4", board.H4, bbFrom(board.H3, board.G3, board.G4, board.G5, board.H5)},
		{"center d4", board.D4, bbFrom(
			board.C3, board.C4, board.C5,
			board.D3, board.D5,
			board.E3, board.E4, board.E5,
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
		sq    board.Square
		color board.Color
		want  board.Bitboard
	}{
		{"white corner a1", board.A1, board.White, bbFrom(board.B2)},
		{"white corner h1", board.H1, board.White, bbFrom(board.G2)},
		{"white edge a4", board.A4, board.White, bbFrom(board.B5)},
		{"white edge h4", board.H4, board.White, bbFrom(board.G5)},
		{"white center d4", board.D4, board.White, bbFrom(board.C5, board.E5)},

		{"black corner a8", board.A8, board.Black, bbFrom(board.B7)},
		{"black corner h8", board.H8, board.Black, bbFrom(board.G7)},
		{"black edge a5", board.A5, board.Black, bbFrom(board.B4)},
		{"black edge h5", board.H5, board.Black, bbFrom(board.G4)},
		{"black center d5", board.D5, board.Black, bbFrom(board.C4, board.E4)},
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
