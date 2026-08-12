package movegen

import (
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func TestIsSquareAttacked(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		sq      board.Square
		byColor board.Color
		want    bool
	}{
		// pawns
		{"white pawn attacks diagonally forward", "4k3/8/8/8/4P3/8/8/4K3 w - - 0 1", board.D5, board.White, true},
		{"white pawn attacks other diagonal", "4k3/8/8/8/4P3/8/8/4K3 w - - 0 1", board.F5, board.White, true},
		{"white pawn does not attack square directly ahead", "4k3/8/8/8/4P3/8/8/4K3 w - - 0 1", board.E5, board.White, false},
		{"black pawn attacks diagonally backward from white's view", "4k3/8/8/4p3/8/8/8/4K3 b - - 0 1", board.D4, board.Black, true},
		{"black pawn does not attack white's diagonal direction", "4k3/8/8/4p3/8/8/8/4K3 b - - 0 1", board.D6, board.Black, false},
		{"black pawn only attacks one diagonal on edge of board", "4k3/8/8/p7/8/8/8/4K3 b - - 0 1", board.H3, board.Black, false},

		// knight
		{"knight attacks standard L-shape", "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", board.F5, board.White, true},
		{"knight does not attack adjacent square", "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", board.D5, board.White, false},
		{"knight on corner a1 attacks only 2 squares", "4k3/8/8/8/8/8/8/N3K3 w - - 0 1", board.B3, board.White, true},
		{"knight on corner a1 does not attack far squares", "4k3/8/8/8/8/8/8/N3K3 w - - 0 1", board.H8, board.White, false},
		{"knight does not wrap across board edge", "4k3/8/8/8/8/8/8/N3K3 w - - 0 1", board.C1, board.White, false},

		// king
		{"king attacks adjacent square", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", board.E2, board.White, true},
		{"king attacks diagonal adjacent square", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", board.D2, board.White, true},
		{"king does not attack two squares away", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", board.E3, board.White, false},
		{"king on corner does not wrap", "3k4/8/8/8/8/8/8/K7 w - - 0 1", board.H1, board.White, false},

		// sliders
		{"rook attacks along open file", "4k3/8/8/4R3/8/8/8/4K3 w - - 0 1", board.E1, board.White, true},
		{"rook attack blocked by own king", "4k3/8/8/4R3/8/8/8/4K3 w - - 0 1", board.D8, board.White, false},
		{"rook attacks along open rank", "4k3/8/8/R3K3/8/8/8/8 w - - 0 1", board.H5, board.White, true},
		{"rook does not attack diagonally", "4k3/8/8/4R3/8/8/8/4K3 w - - 0 1", board.F6, board.White, false},
		{"bishop attacks along open diagonal", "4k3/8/8/8/8/2B5/8/4K3 w - - 0 1", board.H8, board.White, true},
		{"bishop attack blocked by intervening piece", "4k3/6p1/8/8/8/2B5/8/4K3 w - - 0 1", board.H8, board.White, false},
		{"bishop does not attack straight line", "4k3/8/8/8/8/2B5/8/4K3 w - - 0 1", board.C8, board.White, false},
		{"queen attacks diagonally like a bishop", "4k3/8/8/8/8/2Q5/8/4K3 w - - 0 1", board.H8, board.White, true},
		{"queen attacks straight like a rook", "4k3/2Q5/8/8/8/8/8/4K3 w - - 0 1", board.C1, board.White, true},

		// multiple attackers and zero attackers
		{"square attacked by two different pieces still true", "4k3/8/8/4R3/2N5/8/8/4K3 w - - 0 1", board.E4, board.White, true},
		{"square with no attackers at all", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", board.D4, board.White, false},
		{"square occupied by the attacker itself is not being attacked", "4k3/8/8/8/4N3/8/8/4K3 w - - 0 1", board.E4, board.White, false},

		// color filtering - same board, opposite color queried should differ
		{"black piece does not count when querying white attackers", "4k3/8/8/3n4/8/8/8/4K3 w - - 0 1", board.F4, board.White, false},
		{"black knight correctly attacks when querying black", "4k3/8/8/3n4/8/8/8/4K3 w - - 0 1", board.F4, board.Black, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad test FEN: %v", err)
			}
			got := IsSquareAttacked(pos, tt.sq, tt.byColor)
			if got != tt.want {
				t.Errorf("IsSquareAttacked(%s, %s) = %v, want %v", tt.sq, tt.byColor, got, tt.want)
			}
		})
	}
}
