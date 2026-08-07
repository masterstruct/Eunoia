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
