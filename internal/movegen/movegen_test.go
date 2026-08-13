package movegen

import (
	"slices"
	"testing"

	"github.com/masterstruct/Eunoia/internal/board"
)

func TestIsSquareAttacked(t *testing.T) {
	for _, tt := range isSquareAttackedCases {
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

func TestIsSquareAttacked_NoColor(t *testing.T) {
	pos, err := board.ParseFEN("4k3/8/8/8/8/8/8/4K3 w - - 0 1")
	if err != nil {
		t.Fatal(err)
	}

	if IsSquareAttacked(pos, board.E4, board.NoColor) {
		t.Fatal("expected false for NoColor")
	}
}

func BenchmarkIsSquareAttacked(b *testing.B) {
	positions := make([]board.Position, len(isSquareAttackedCases))
	for i, tc := range isSquareAttackedCases {
		pos, err := board.ParseFEN(tc.fen)
		if err != nil {
			b.Fatalf("bad bench FEN: %v", err)
		}
		positions[i] = pos
	}

	for i := 0; b.Loop(); i++ {
		tc := isSquareAttackedCases[i%len(isSquareAttackedCases)]
		pos := positions[i%len(positions)]
		_ = IsSquareAttacked(pos, tc.sq, tc.byColor)
	}
}

var isSquareAttackedCases = []struct {
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
	{"rook attacks along open file", "4k3/8/8/4R3/8/8/8/4K3 w - - 0 1", board.A5, board.White, true},
	{"rook attack blocked by own king", "4k3/8/4K3/4R3/8/8/8/8 w - - 0 1", board.E8, board.White, false},
	{"rook attacks along open rank", "4k3/8/8/R7/4K3/8/8/8 w - - 0 1", board.H5, board.White, true},
	{"rook does not attack diagonally", "4k3/8/8/4R3/8/8/8/4K3 w - - 0 1", board.F6, board.White, false},
	{"bishop attacks along open diagonal", "4k3/8/8/8/8/2B5/8/4K3 w - - 0 1", board.H8, board.White, true},
	{"bishop attack blocked by intervening piece", "4k3/6p1/8/8/8/2B5/8/4K3 w - - 0 1", board.H8, board.White, false},
	{"bishop does not attack straight line", "4k3/8/8/8/8/2B5/8/4K3 w - - 0 1", board.C8, board.White, false},
	{"queen attacks diagonally like a bishop", "4k3/8/8/8/8/2Q5/8/4K3 w - - 0 1", board.H8, board.White, true},
	{"queen attacks straight like a rook", "4k3/2Q5/8/8/8/8/8/4K3 w - - 0 1", board.C1, board.White, true},

	// multiple attackers and zero attackers
	{"square attacked by two different pieces", "4k3/8/8/4R3/2N5/8/8/4K3 w - - 0 1", board.E3, board.White, true},
	{"square with no attackers at all", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", board.D4, board.White, false},
	{"square occupied by the attacker itself is not being attacked", "4k3/8/8/8/4N3/8/8/4K3 w - - 0 1", board.E4, board.White, false},

	// color filtering - same board, opposite color queried should differ
	{"black piece does not count when querying white attackers", "4k3/8/8/3n4/8/8/8/4K3 w - - 0 1", board.C3, board.White, false},
	{"black knight correctly attacks when querying black", "4k3/8/8/3n4/8/8/8/4K3 w - - 0 1", board.C3, board.Black, true},

	// attacking occupied squares
	{"white rook attacks black piece", "8/2k5/8/1n2R3/8/8/8/4K3 w - - 0 1", board.B5, board.White, true},
	{"black rook attacks white piece", "8/2k5/8/5r2/8/8/5B2/4K3 b - - 0 1", board.F2, board.Black, true},
	{"white rook attacks (defends) white piece", "8/2k5/8/1Q2R3/8/8/8/4K3 w - - 0 1", board.B5, board.White, true},
	{"black rook attacks (defends) black piece", "8/2k5/8/5r2/8/8/5n2/4K3 b - - 0 1", board.F2, board.Black, true},
	{"white bishop attacks black piece", "8/8/6k1/6n1/8/4B3/1K6/8 w - - 0 1", board.G5, board.White, true},
	{"black queen attacks white piece", "8/8/6k1/3q4/8/3N4/1K6/8 b - - 0 1", board.D3, board.Black, true},
}

func TestInCheck(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want bool
	}{
		{
			name: "white in check by rook",
			fen:  "4k3/8/8/8/8/8/4r3/4K3 w - - 0 1",
			want: true,
		},
		{
			name: "black in check by bishop",
			fen:  "8/5k2/8/8/2B5/5K2/8/8 b - - 0 1",
			want: true,
		},
		{
			name: "white not in check",
			fen:  "4k3/8/3r4/1b6/8/5q2/1n6/4K3 w - - 0 1",
			want: false,
		},
		{
			name: "black not in check",
			fen:  "8/3k4/8/8/8/8/4R3/4K3 b - - 0 1",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := board.ParseFEN(tt.fen)

			got := InCheck(pos)
			if got != tt.want {
				t.Fatalf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestGenKnightMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []board.Square
	}{
		{
			name: "starting position for white has 4 knight moves",
			fen:  board.StartingFEN,
			to:   []board.Square{board.A3, board.C3, board.F3, board.H3},
		},
		{
			name: "starting position for black has 4 knight moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []board.Square{board.A6, board.C6, board.F6, board.H6},
		},
		{
			name: "kiwipete for white has 11 knight moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
			to: []board.Square{board.B1, board.D1, board.A4, board.B5,
				board.D3, board.C4, board.G4, board.C6, board.G6, board.D7, board.F7},
		},
		{
			name: "kiwipete for black has 10 knight moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []board.Square{board.A4, board.C4, board.D5, board.C8,
				board.E4, board.G4, board.D5, board.H5, board.G8, board.H7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := board.ParseFEN(tt.fen)

			movelist := genKnightMoves(pos)
			for _, move := range movelist {
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected knight move: %v\n%v", move, pos)
				}
			}
			if len(movelist) != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos)
			}
		})
	}
}

func TestGenBishopMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []board.Square
	}{
		{
			name: "starting position for white has 0 bishop moves",
			fen:  board.StartingFEN,
			to:   []board.Square{},
		},
		{
			name: "starting position for black has 0 bishop moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []board.Square{},
		},
		{
			name: "kiwipete for white has 11 bishop moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
			to: []board.Square{board.C1, board.E3, board.F4, board.G5, board.H6,
				board.A6, board.B5, board.C4, board.D3, board.D1, board.F1},
		},
		{
			name: "kiwipete for black has 8 bishop moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []board.Square{board.B5, board.B7, board.C4,
				board.C8, board.D3, board.E2, board.F8, board.H6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := board.ParseFEN(tt.fen)

			movelist := genBishopMoves(pos)
			for _, move := range movelist {
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected bishop move: %v\n%v", move, pos)
				}
			}
			if len(movelist) != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos)
			}
		})
	}
}

func TestGenRookMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []board.Square
	}{
		{
			name: "starting position for white has 0 rook moves",
			fen:  board.StartingFEN,
			to:   []board.Square{},
		},
		{
			name: "starting position for black has 0 rook moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []board.Square{},
		},
		{
			name: "kiwipete for white has 5 rook moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
			to:   []board.Square{board.B1, board.C1, board.D1, board.F1, board.G1},
		},
		{
			name: "kiwipete for black has 9 rook moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []board.Square{board.B8, board.C8, board.D8,
				board.F8, board.G8, board.H4, board.H5, board.H6, board.H7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := board.ParseFEN(tt.fen)

			movelist := genRookMoves(pos)
			for _, move := range movelist {
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected rook move: %v\n%v", move, pos)
				}
			}
			if len(movelist) != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos)
			}
		})
	}
}

func TestGenQueenMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []board.Square
	}{
		{
			name: "starting position for white has 0 queen moves",
			fen:  board.StartingFEN,
			to:   []board.Square{},
		},
		{
			name: "starting position for black has 0 queen moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []board.Square{},
		},
		{
			name: "kiwipete for white has 9 queen moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
			to: []board.Square{board.D3, board.E3, board.F4, board.F5,
				board.F6, board.G3, board.G4, board.H3, board.H5},
		},
		{
			name: "kiwipete for black has 4 queen moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to:   []board.Square{board.C5, board.D6, board.D8, board.F8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := board.ParseFEN(tt.fen)

			movelist := genQueenMoves(pos)
			for _, move := range movelist {
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected queen move: %v\n%v", move, pos)
				}
			}
			if len(movelist) != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos)
			}
		})
	}
}

func TestGenPawnMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []board.Square
	}{
		{
			name: "starting position for white has 16 pawn moves",
			fen:  board.StartingFEN,
			to: []board.Square{
				board.A3, board.A4, board.B3, board.B4, board.C3, board.C4, board.D3, board.D4,
				board.E3, board.E4, board.F3, board.F4, board.G3, board.G4, board.H3, board.H4,
			},
		},
		{
			name: "starting position for black has 16 pawn moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to: []board.Square{
				board.A6, board.A5, board.B6, board.B5, board.C6, board.C5, board.D6, board.D5,
				board.E6, board.E5, board.F6, board.F5, board.G6, board.G5, board.H6, board.H5,
			},
		},
		{
			name: "kiwipete for white has 8 pawn moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
			to: []board.Square{board.A3, board.A4, board.B3,
				board.D6, board.E6, board.G3, board.G4, board.H3},
		},
		{
			name: "kiwipete for black has 8 pawn moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []board.Square{board.C6, board.C5, board.D6,
				board.D5, board.G5, board.B3, board.C3, board.G2},
		},
		{
			name: "single and double push blocked by piece",
			fen:  "4k3/8/8/8/8/2pp4/3P4/4K3 w - - 0 1",
			to:   []board.Square{board.C3},
		},
		{
			name: "double push blocked by piece",
			fen:  "4k3/8/8/8/3p4/8/3P4/4K3 w - - 0 1",
			to:   []board.Square{board.D3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			movelist := genPawnMoves(pos)
			for _, move := range movelist {
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected pawn move: %v\n%v", move, pos)
				}
			}
			if len(movelist) != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos)
			}
		})
	}
}

func TestGenPawnMoves_EnPassant(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []board.Square
	}{
		{
			name: "white captures en passant",
			fen:  "4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
			to:   []board.Square{board.E6, board.D6},
		},
		{
			name: "black captures en passant",
			fen:  "4k3/8/8/8/3pP3/8/8/4K3 b - e3 0 1",
			to:   []board.Square{board.D3, board.E3},
		},
		{
			name: "no en passant square set",
			fen:  "4k3/8/8/3pP3/8/8/8/4K3 w - - 0 1",
			to:   []board.Square{board.E6},
		},
		{
			name: "two pawns can capture en passant",
			fen:  "4k3/8/8/1PpP4/8/8/8/4K3 w - c6 0 1",
			to:   []board.Square{board.B6, board.C6, board.D6, board.C6},
		},
		{
			name: "en passant pin: pseudolegal movegen still creates capture",
			fen:  "8/2p5/3p4/KP5r/1R3pPk/8/8/8 b - g3 0 1",
			to:   []board.Square{board.F3, board.G3, board.D5, board.C6, board.C5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			movelist := genPawnMoves(pos)
			for _, move := range movelist {
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected pawn move: %v\n%v", move, pos)
				}
			}
			if len(movelist) != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos)
			}
		})
	}
}

func TestGenPawnMoves_Promotion(t *testing.T) {
	tests := []struct {
		name      string
		fen       string
		wantMoves []board.Move
	}{
		{
			name: "white promotions and capture promotions",
			fen:  "3r4/4P3/8/8/8/8/8/4K2k w - - 0 1",
			wantMoves: []board.Move{
				board.NewCapturePromo(board.E7, board.D8, board.Knight),
				board.NewCapturePromo(board.E7, board.D8, board.Bishop),
				board.NewCapturePromo(board.E7, board.D8, board.Rook),
				board.NewCapturePromo(board.E7, board.D8, board.Queen),
				board.NewPromo(board.E7, board.E8, board.Knight),
				board.NewPromo(board.E7, board.E8, board.Bishop),
				board.NewPromo(board.E7, board.E8, board.Rook),
				board.NewPromo(board.E7, board.E8, board.Queen),
			},
		},
		{
			name: "white capture promotions only",
			fen:  "3rn3/4P3/8/8/8/8/8/4K2k w - - 0 1",
			wantMoves: []board.Move{
				board.NewCapturePromo(board.E7, board.D8, board.Knight),
				board.NewCapturePromo(board.E7, board.D8, board.Bishop),
				board.NewCapturePromo(board.E7, board.D8, board.Rook),
				board.NewCapturePromo(board.E7, board.D8, board.Queen),
			},
		},
		{
			name: "white promotions only",
			fen:  "8/4P3/8/8/8/8/8/4K2k w - - 0 1",
			wantMoves: []board.Move{
				board.NewPromo(board.E7, board.E8, board.Knight),
				board.NewPromo(board.E7, board.E8, board.Bishop),
				board.NewPromo(board.E7, board.E8, board.Rook),
				board.NewPromo(board.E7, board.E8, board.Queen),
			},
		},
		{
			name: "black promotions and capture promotions",
			fen:  "8/8/8/8/7k/1K6/4p3/3R4 b - - 0 1",
			wantMoves: []board.Move{
				board.NewCapturePromo(board.E2, board.D1, board.Knight),
				board.NewCapturePromo(board.E2, board.D1, board.Bishop),
				board.NewCapturePromo(board.E2, board.D1, board.Rook),
				board.NewCapturePromo(board.E2, board.D1, board.Queen),
				board.NewPromo(board.E2, board.E1, board.Knight),
				board.NewPromo(board.E2, board.E1, board.Bishop),
				board.NewPromo(board.E2, board.E1, board.Rook),
				board.NewPromo(board.E2, board.E1, board.Queen),
			},
		},
		{
			name: "black capture promotions only",
			fen:  "8/8/8/8/7k/1K6/4p3/3RN3 b - - 0 1",
			wantMoves: []board.Move{
				board.NewCapturePromo(board.E2, board.D1, board.Knight),
				board.NewCapturePromo(board.E2, board.D1, board.Bishop),
				board.NewCapturePromo(board.E2, board.D1, board.Rook),
				board.NewCapturePromo(board.E2, board.D1, board.Queen),
			},
		},
		{
			name: "black promotions only",
			fen:  "8/8/8/8/7k/1K6/4p3/8 b - - 0 1",
			wantMoves: []board.Move{
				board.NewPromo(board.E2, board.E1, board.Knight),
				board.NewPromo(board.E2, board.E1, board.Bishop),
				board.NewPromo(board.E2, board.E1, board.Rook),
				board.NewPromo(board.E2, board.E1, board.Queen),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			movelist := genPawnMoves(pos)
			for _, move := range movelist {
				if !slices.Contains(tt.wantMoves, move) {
					t.Fatalf("unexpected pawn move: %v\n%v", move, pos)
				}
			}
			if len(movelist) != len(tt.wantMoves) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.wantMoves, movelist, pos)
			}
		})
	}
}

// TODO: figure out castling pseudolegal checks

// func TestGenKingMoves(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		fen  string
// 		to   []board.Square
// 	}{
// 		{
// 			name: "starting position for white has 0 king moves",
// 			fen:  board.StartingFEN,
// 			to:   []board.Square{},
// 		},
// 		{
// 			name: "center of board king has 8 moves",
// 			fen:  "8/8/2k5/8/8/4K3/8/8 w - - 0 1",
// 			to: []board.Square{board.D2, board.D3, board.D4,
// 				board.E4, board.F4, board.F3, board.F2, board.E2},
// 		},
// 		{
// 			name: "king attackers restrict movement",
// 			fen:  "8/8/8/3k4/8/4K3/8/6n1 w - - 0 1",
// 			to:   []board.Square{board.D2, board.D3, board.F4, board.F2},
// 		},
// 		{
// 			name: "kiwipete for white has 4 king moves",
// 			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
// 			to:   []board.Square{board.C1, board.D1, board.F1, board.G1},
// 		},
// 		{
// 			name: "kiwipete for black has 4 king moves",
// 			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
// 			to:   []board.Square{board.C8, board.D8, board.F8, board.G8},
// 		},
// 		{
// 			name: "white king can castle kingside",
// 			fen:  "8/8/8/3k4/8/8/8/4K2R w K - 0 1",
// 			to: []board.Square{board.D1, board.D2, board.E2,
// 				board.F2, board.F1, board.H1},
// 		},
// 		{
// 			name: "white king cannot castle kingside through check",
// 			fen:  "k7/8/8/8/2b5/8/8/4K2R w K - 0 1",
// 			to:   []board.Square{board.D1, board.D2, board.F2},
// 		},
// 		{
// 			name: "white king can castle queenside",
// 			fen:  "8/8/8/3k4/8/8/8/R3K3 w Q - 0 1",
// 			to: []board.Square{board.A1, board.D1, board.D2,
// 				board.E2, board.F2, board.F1},
// 		},
// 		{
// 			name: "white king cannot castle queenside through check",
// 			fen:  "8/8/8/2k5/6b1/8/8/R3K3 w Q - 0 1",
// 			to:   []board.Square{board.D2, board.F2, board.F1},
// 		},
// 		{
// 			name: "white king can castle both ways",
// 			fen:  "8/8/8/3k4/8/8/8/R3K2R w KQ - 0 1",
// 			to: []board.Square{board.A1, board.D1, board.D2,
// 				board.E2, board.F2, board.F1, board.H1},
// 		},
// 		{
// 			name: "black king can castle kingside",
// 			fen:  "4k2r/8/8/3K4/8/8/8/8 b k - 0 1",
// 			to: []board.Square{board.D8, board.D7, board.E7,
// 				board.F7, board.F8, board.H8},
// 		},
// 		{
// 			name: "black king cannot castle kingside through check",
// 			fen:  "r3k2r/8/3B4/3K4/8/8/8/8 b kq - 0 1",
// 			to:   []board.Square{board.A8, board.D8, board.D7, board.F7},
// 		},
// 		{
// 			name: "black king can castle queenside",
// 			fen:  "r3k3/8/8/3K4/8/8/8/8 b q - 0 1",
// 			to: []board.Square{board.A8, board.D8, board.D7,
// 				board.E7, board.F7, board.F8},
// 		},
// 		{
// 			name: "black king cannot castle queenside through check",
// 			fen:  "r3k2r/8/8/3K2B1/8/8/8/8 b kq - 0 1",
// 			to:   []board.Square{board.D7, board.F7, board.F8, board.H8},
// 		},
// 		{
// 			name: "black king can castle both ways",
// 			fen:  "r3k2r/8/8/3K4/8/8/8/8 b kq - 0 1",
// 			to: []board.Square{board.A8, board.D8, board.D7,
// 				board.E7, board.F7, board.F8, board.H8},
// 		},
// 		{
// 			name: "king cannot castle to attacked square",
// 			fen:  "8/8/8/3k4/8/7n/8/4K2R w K - 0 1",
// 			to:   []board.Square{board.D1, board.D2, board.E2, board.F1},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			pos, _ := board.ParseFEN(tt.fen)

// 			movelist := genKingMoves(pos)
// 			for _, move := range movelist {
// 				if !slices.Contains(tt.to, move.To()) {
// 					t.Fatalf("unexpected knight move: %v\n%v", move, pos)
// 				}
// 			}
// 			if len(movelist) != len(tt.to) {
// 				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos)
// 			}
// 		})
// 	}
// }
