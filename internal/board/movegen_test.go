package board

import (
	"slices"
	"testing"
)

func TestIsSquareAttacked(t *testing.T) {
	for _, tt := range isSquareAttackedCases {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad test FEN: %v", err)
			}
			got := IsSquareAttacked(&pos, tt.sq, tt.byColor)
			if got != tt.want {
				t.Errorf("IsSquareAttacked(%s, %s) = %v, want %v", tt.sq, tt.byColor, got, tt.want)
			}
		})
	}
}

func TestIsSquareAttacked_NoColor(t *testing.T) {
	pos, err := ParseFEN("4k3/8/8/8/8/8/8/4K3 w - - 0 1")
	if err != nil {
		t.Fatal(err)
	}

	if IsSquareAttacked(&pos, E4, NoColor) {
		t.Fatal("expected false for NoColor")
	}
}

func BenchmarkIsSquareAttacked(b *testing.B) {
	positions := make([]Position, len(isSquareAttackedCases))
	for i, tc := range isSquareAttackedCases {
		pos, err := ParseFEN(tc.fen)
		if err != nil {
			b.Fatalf("bad bench FEN: %v", err)
		}
		positions[i] = pos
	}

	for i := 0; b.Loop(); i++ {
		tc := isSquareAttackedCases[i%len(isSquareAttackedCases)]
		pos := positions[i%len(positions)]
		_ = IsSquareAttacked(&pos, tc.sq, tc.byColor)
	}
}

var isSquareAttackedCases = []struct {
	name    string
	fen     string
	sq      Square
	byColor Color
	want    bool
}{
	{"white pawn attacking black knight", "k7/8/8/2n5/3P4/4K3/8/8 b - - 0 1", C5, White, true},
	// pawns
	{"white pawn attacks diagonally forward", "4k3/8/8/8/4P3/8/8/4K3 w - - 0 1", D5, White, true},
	{"white pawn attacks other diagonal", "4k3/8/8/8/4P3/8/8/4K3 w - - 0 1", F5, White, true},
	{"white pawn does not attack square directly ahead", "4k3/8/8/8/4P3/8/8/4K3 w - - 0 1", E5, White, false},
	{"black pawn attacks diagonally backward from white's view", "4k3/8/8/4p3/8/8/8/4K3 b - - 0 1", D4, Black, true},
	{"black pawn does not attack white's diagonal direction", "4k3/8/8/4p3/8/8/8/4K3 b - - 0 1", D6, Black, false},
	{"black pawn only attacks one diagonal on edge of board", "4k3/8/8/p7/8/8/8/4K3 b - - 0 1", H3, Black, false},

	// knight
	{"knight attacks standard L-shape", "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", F5, White, true},
	{"knight does not attack adjacent square", "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", D5, White, false},
	{"knight on corner a1 attacks only 2 squares", "4k3/8/8/8/8/8/8/N3K3 w - - 0 1", B3, White, true},
	{"knight on corner a1 does not attack far squares", "4k3/8/8/8/8/8/8/N3K3 w - - 0 1", H8, White, false},
	{"knight does not wrap across board edge", "4k3/8/8/8/8/8/8/N3K3 w - - 0 1", C1, White, false},

	// king
	{"king attacks adjacent square", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", E2, White, true},
	{"king attacks diagonal adjacent square", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", D2, White, true},
	{"king does not attack two squares away", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", E3, White, false},
	{"king on corner does not wrap", "3k4/8/8/8/8/8/8/K7 w - - 0 1", H1, White, false},

	// sliders
	{"rook attacks along open file", "4k3/8/8/4R3/8/8/8/4K3 w - - 0 1", A5, White, true},
	{"rook attack blocked by own king", "4k3/8/4K3/4R3/8/8/8/8 w - - 0 1", E8, White, false},
	{"rook attacks along open rank", "4k3/8/8/R7/4K3/8/8/8 w - - 0 1", H5, White, true},
	{"rook does not attack diagonally", "4k3/8/8/4R3/8/8/8/4K3 w - - 0 1", F6, White, false},
	{"bishop attacks along open diagonal", "4k3/8/8/8/8/2B5/8/4K3 w - - 0 1", H8, White, true},
	{"bishop attack blocked by intervening piece", "4k3/6p1/8/8/8/2B5/8/4K3 w - - 0 1", H8, White, false},
	{"bishop does not attack straight line", "4k3/8/8/8/8/2B5/8/4K3 w - - 0 1", C8, White, false},
	{"queen attacks diagonally like a bishop", "4k3/8/8/8/8/2Q5/8/4K3 w - - 0 1", H8, White, true},
	{"queen attacks straight like a rook", "4k3/2Q5/8/8/8/8/8/4K3 w - - 0 1", C1, White, true},

	// multiple attackers and zero attackers
	{"square attacked by two different pieces", "4k3/8/8/4R3/2N5/8/8/4K3 w - - 0 1", E3, White, true},
	{"square with no attackers at all", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", D4, White, false},
	{"square occupied by the attacker itself is not being attacked", "4k3/8/8/8/4N3/8/8/4K3 w - - 0 1", E4, White, false},

	// color filtering - same board, opposite color queried should differ
	{"black piece does not count when querying white attackers", "4k3/8/8/3n4/8/8/8/4K3 w - - 0 1", C3, White, false},
	{"black knight correctly attacks when querying black", "4k3/8/8/3n4/8/8/8/4K3 w - - 0 1", C3, Black, true},

	// attacking occupied squares
	{"white rook attacks black piece", "8/2k5/8/1n2R3/8/8/8/4K3 w - - 0 1", B5, White, true},
	{"black rook attacks white piece", "8/2k5/8/5r2/8/8/5B2/4K3 b - - 0 1", F2, Black, true},
	{"white rook attacks (defends) white piece", "8/2k5/8/1Q2R3/8/8/8/4K3 w - - 0 1", B5, White, true},
	{"black rook attacks (defends) black piece", "8/2k5/8/5r2/8/8/5n2/4K3 b - - 0 1", F2, Black, true},
	{"white bishop attacks black piece", "8/8/6k1/6n1/8/4B3/1K6/8 w - - 0 1", G5, White, true},
	{"black queen attacks white piece", "8/8/6k1/3q4/8/3N4/1K6/8 b - - 0 1", D3, Black, true},
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
			name: "white king not in check in kiwipete position",
			fen:  KiwipeteFEN,
			want: false,
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
			pos, _ := ParseFEN(tt.fen)

			got := InCheck(&pos, pos.SideToMove)
			if got != tt.want {
				t.Fatalf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func BenchmarkInCheck(b *testing.B) {
	fens := []string{
		"4k3/8/8/8/8/8/4r3/4K3 w - - 0 1",
		KiwipeteFEN,
		"8/5k2/8/8/2B5/5K2/8/8 b - - 0 1",
		"4k3/8/3r4/1b6/8/5q2/1n6/4K3 w - - 0 1",
		"8/3k4/8/8/8/8/4R3/4K3 b - - 0 1",
	}

	positions := make([]Position, len(fens))
	for i, fen := range fens {
		pos, _ := ParseFEN(fen)
		positions[i] = pos
	}

	for i := 0; b.Loop(); i++ {
		pos := &positions[i%len(positions)]
		_ = InCheck(pos, pos.SideToMove)
	}
}

func TestGenKnightMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []Square
	}{
		{
			name: "starting position for white has 4 knight moves",
			fen:  StartingFEN,
			to:   []Square{A3, C3, F3, H3},
		},
		{
			name: "starting position for black has 4 knight moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []Square{A6, C6, F6, H6},
		},
		{
			name: "kiwipete for white has 11 knight moves",
			fen:  KiwipeteFEN,
			to: []Square{B1, D1, A4, B5,
				D3, C4, G4, C6, G6, D7, F7},
		},
		{
			name: "kiwipete for black has 10 knight moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []Square{A4, C4, D5, C8,
				E4, G4, D5, H5, G8, H7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := ParseFEN(tt.fen)

			var movelist Movelist
			GenKnightMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected knight move: %v\n%v", move, pos)
				}
			}
			if movelist.Len != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos.String())
			}
		})
	}
}

func TestGenBishopMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []Square
	}{
		{
			name: "starting position for white has 0 bishop moves",
			fen:  StartingFEN,
			to:   []Square{},
		},
		{
			name: "starting position for black has 0 bishop moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []Square{},
		},
		{
			name: "kiwipete for white has 11 bishop moves",
			fen:  KiwipeteFEN,
			to: []Square{C1, E3, F4, G5, H6,
				A6, B5, C4, D3, D1, F1},
		},
		{
			name: "kiwipete for black has 8 bishop moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []Square{B5, B7, C4,
				C8, D3, E2, F8, H6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := ParseFEN(tt.fen)

			var movelist Movelist
			GenBishopMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected bishop move: %v\n%v", move, pos)
				}
			}
			if movelist.Len != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos.String())
			}
		})
	}
}

func TestGenRookMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []Square
	}{
		{
			name: "starting position for white has 0 rook moves",
			fen:  StartingFEN,
			to:   []Square{},
		},
		{
			name: "starting position for black has 0 rook moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []Square{},
		},
		{
			name: "kiwipete for white has 5 rook moves",
			fen:  KiwipeteFEN,
			to:   []Square{B1, C1, D1, F1, G1},
		},
		{
			name: "kiwipete for black has 9 rook moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []Square{B8, C8, D8,
				F8, G8, H4, H5, H6, H7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := ParseFEN(tt.fen)

			var movelist Movelist
			GenRookMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected rook move: %v\n%v", move, pos)
				}
			}
			if movelist.Len != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos.String())
			}
		})
	}
}

func TestGenQueenMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []Square
	}{
		{
			name: "starting position for white has 0 queen moves",
			fen:  StartingFEN,
			to:   []Square{},
		},
		{
			name: "starting position for black has 0 queen moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to:   []Square{},
		},
		{
			name: "kiwipete for white has 9 queen moves",
			fen:  KiwipeteFEN,
			to: []Square{D3, E3, F4, F5,
				F6, G3, G4, H3, H5},
		},
		{
			name: "kiwipete for black has 4 queen moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to:   []Square{C5, D6, D8, F8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := ParseFEN(tt.fen)

			var movelist Movelist
			GenQueenMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected queen move: %v\n%v", move, pos)
				}
			}
			if movelist.Len != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos.String())
			}
		})
	}
}

func TestGenPawnMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []Square
	}{
		{
			name: "starting position for white has 16 pawn moves",
			fen:  StartingFEN,
			to: []Square{
				A3, A4, B3, B4, C3, C4, D3, D4,
				E3, E4, F3, F4, G3, G4, H3, H4,
			},
		},
		{
			name: "starting position for black has 16 pawn moves",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			to: []Square{
				A6, A5, B6, B5, C6, C5, D6, D5,
				E6, E5, F6, F5, G6, G5, H6, H5,
			},
		},
		{
			name: "kiwipete for white has 8 pawn moves",
			fen:  KiwipeteFEN,
			to: []Square{A3, A4, B3,
				D6, E6, G3, G4, H3},
		},
		{
			name: "kiwipete for black has 8 pawn moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to: []Square{C6, C5, D6,
				D5, G5, B3, C3, G2},
		},
		{
			name: "single and double push blocked by piece",
			fen:  "4k3/8/8/8/8/2pp4/3P4/4K3 w - - 0 1",
			to:   []Square{C3},
		},
		{
			name: "double push blocked by piece",
			fen:  "4k3/8/8/8/3p4/8/3P4/4K3 w - - 0 1",
			to:   []Square{D3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			var movelist Movelist
			GenPawnMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected pawn move: %v\n%v", move, pos.String())
				}
			}
			if movelist.Len != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos.String())
			}
		})
	}
}

func TestGenPawnMoves_EnPassant(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []Square
	}{
		{
			name: "white captures en passant",
			fen:  "4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
			to:   []Square{E6, D6},
		},
		{
			name: "black captures en passant",
			fen:  "4k3/8/8/8/3pP3/8/8/4K3 b - e3 0 1",
			to:   []Square{D3, E3},
		},
		{
			name: "no en passant square set",
			fen:  "4k3/8/8/3pP3/8/8/8/4K3 w - - 0 1",
			to:   []Square{E6},
		},
		{
			name: "two pawns can capture en passant",
			fen:  "4k3/8/8/1PpP4/8/8/8/4K3 w - c6 0 1",
			to:   []Square{B6, C6, D6, C6},
		},
		{
			name: "en passant pin: pseudolegal movegen still creates capture",
			fen:  "8/2p5/3p4/KP5r/1R3pPk/8/8/8 b - g3 0 1",
			to:   []Square{F3, G3, D5, C6, C5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			var movelist Movelist
			GenPawnMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected pawn move: %v\n%v", move, pos.String())
				}
			}
			if movelist.Len != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos.String())
			}
		})
	}
}

func TestGenPawnMoves_Promotion(t *testing.T) {
	tests := []struct {
		name      string
		fen       string
		wantMoves []Move
	}{
		{
			name: "white promotions and capture promotions",
			fen:  "3r4/4P3/8/8/8/8/8/4K2k w - - 0 1",
			wantMoves: []Move{
				NewCapturePromo(E7, D8, Knight),
				NewCapturePromo(E7, D8, Bishop),
				NewCapturePromo(E7, D8, Rook),
				NewCapturePromo(E7, D8, Queen),
				NewPromo(E7, E8, Knight),
				NewPromo(E7, E8, Bishop),
				NewPromo(E7, E8, Rook),
				NewPromo(E7, E8, Queen),
			},
		},
		{
			name: "white capture promotions only",
			fen:  "3rn3/4P3/8/8/8/8/8/4K2k w - - 0 1",
			wantMoves: []Move{
				NewCapturePromo(E7, D8, Knight),
				NewCapturePromo(E7, D8, Bishop),
				NewCapturePromo(E7, D8, Rook),
				NewCapturePromo(E7, D8, Queen),
			},
		},
		{
			name: "white promotions only",
			fen:  "8/4P3/8/8/8/8/8/4K2k w - - 0 1",
			wantMoves: []Move{
				NewPromo(E7, E8, Knight),
				NewPromo(E7, E8, Bishop),
				NewPromo(E7, E8, Rook),
				NewPromo(E7, E8, Queen),
			},
		},
		{
			name: "black promotions and capture promotions",
			fen:  "8/8/8/8/7k/1K6/4p3/3R4 b - - 0 1",
			wantMoves: []Move{
				NewCapturePromo(E2, D1, Knight),
				NewCapturePromo(E2, D1, Bishop),
				NewCapturePromo(E2, D1, Rook),
				NewCapturePromo(E2, D1, Queen),
				NewPromo(E2, E1, Knight),
				NewPromo(E2, E1, Bishop),
				NewPromo(E2, E1, Rook),
				NewPromo(E2, E1, Queen),
			},
		},
		{
			name: "black capture promotions only",
			fen:  "8/8/8/8/7k/1K6/4p3/3RN3 b - - 0 1",
			wantMoves: []Move{
				NewCapturePromo(E2, D1, Knight),
				NewCapturePromo(E2, D1, Bishop),
				NewCapturePromo(E2, D1, Rook),
				NewCapturePromo(E2, D1, Queen),
			},
		},
		{
			name: "black promotions only",
			fen:  "8/8/8/8/7k/1K6/4p3/8 b - - 0 1",
			wantMoves: []Move{
				NewPromo(E2, E1, Knight),
				NewPromo(E2, E1, Bishop),
				NewPromo(E2, E1, Rook),
				NewPromo(E2, E1, Queen),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			var movelist Movelist
			GenPawnMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.wantMoves, move) {
					t.Fatalf("unexpected pawn move: %v\n%v", move, pos)
				}
			}
			if movelist.Len != len(tt.wantMoves) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.wantMoves, movelist, pos.String())
			}
		})
	}
}

func TestGenKingMoves(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		to   []Square
	}{
		{
			name: "starting position for white has 0 king moves",
			fen:  StartingFEN,
			to:   []Square{},
		},
		{
			name: "center of board king has 8 moves",
			fen:  "8/8/2k5/8/8/4K3/8/8 w - - 0 1",
			to: []Square{D2, D3, D4,
				E4, F4, F3, F2, E2},
		},
		{
			name: "kiwipete for white has 4 king moves",
			fen:  KiwipeteFEN,
			to:   []Square{A1, D1, F1, H1},
		},
		{
			name: "kiwipete for black has 4 king moves",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq -",
			to:   []Square{A8, D8, F8, H8},
		},
		{
			name: "white king can castle kingside",
			fen:  "8/8/8/3k4/8/8/8/4K2R w K - 0 1",
			to: []Square{D1, D2, E2,
				F2, F1, H1},
		},
		{
			name: "white king cannot castle kingside through check",
			fen:  "k7/8/8/8/2b5/8/8/4K2R w K - 0 1",
			to: []Square{D1, D2,
				E2, F2, F1},
		},
		{
			name: "white king can castle queenside",
			fen:  "8/8/8/3k4/8/8/8/R3K3 w Q - 0 1",
			to: []Square{A1, D1, D2,
				E2, F2, F1},
		},
		{
			name: "white king cannot castle queenside through check",
			fen:  "8/8/8/2k5/6b1/8/8/R3K3 w Q - 0 1",
			to: []Square{D1, D2,
				E2, F2, F1},
		},
		{
			name: "white king can castle both ways",
			fen:  "8/8/8/3k4/8/8/8/R3K2R w KQ - 0 1",
			to: []Square{A1, D1, D2,
				E2, F2, F1, H1},
		},
		{
			name: "black king can castle kingside",
			fen:  "4k2r/8/8/3K4/8/8/8/8 b k - 0 1",
			to: []Square{D8, D7, E7,
				F7, F8, H8},
		},
		{
			name: "black king cannot castle kingside through check",
			fen:  "r3k2r/8/3B4/3K4/8/8/8/8 b kq - 0 1",
			to: []Square{A8, D8, D7,
				E7, F7, F8},
		},
		{
			name: "black king can castle queenside",
			fen:  "r3k3/8/8/3K4/8/8/8/8 b q - 0 1",
			to: []Square{A8, D8, D7,
				E7, F7, F8},
		},
		{
			name: "black king cannot castle queenside through check",
			fen:  "r3k2r/8/8/3K2B1/8/8/8/8 b kq - 0 1",
			to: []Square{D8, D7, E7,
				F7, F8, H8},
		},
		{
			name: "black king can castle both ways",
			fen:  "r3k2r/8/8/3K4/8/8/8/8 b kq - 0 1",
			to: []Square{A8, D8, D7,
				E7, F7, F8, H8},
		},
		{
			name: "king cannot castle to attacked square",
			fen:  "8/8/8/3k4/8/7n/8/4K2R w K - 0 1",
			to:   []Square{D1, D2, E2, F2, F1},
		},
		{
			name: "king cannot castle if check",
			fen:  "8/8/8/3k4/8/3n4/8/4K2R w K - 0 1",
			to:   []Square{D1, D2, E2, F2, F1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := ParseFEN(tt.fen)

			var movelist Movelist
			GenKingMoves(&pos, &movelist)
			for i := range movelist.Len {
				move := movelist.Moves[i]
				if !slices.Contains(tt.to, move.To()) {
					t.Fatalf("unexpected king move: %v\n%v", move, pos.String())
				}
			}
			if movelist.Len != len(tt.to) {
				t.Fatalf("missing moves: expected %v but got %v\n%v", tt.to, movelist, pos.String())
			}
		})
	}
}

func TestCanCastle(t *testing.T) {
	tests := []struct {
		name   string
		fen    string
		kingSq Square
		rookSq Square
		toFile File
		dir    File
		want   bool
	}{
		{
			name:   "clear kingside path",
			fen:    "4k3/8/8/8/8/8/8/4K2R w K - 0 1",
			kingSq: E1,
			rookSq: H1,
			toFile: FileG,
			dir:    1,
			want:   true,
		},
		{
			name:   "castling path blocked",
			fen:    "8/8/8/8/2k5/8/8/4KN1R w K - 0 1",
			kingSq: E1,
			rookSq: H1,
			toFile: FileG,
			dir:    1,
			want:   false,
		},
		{
			name:   "destination square occupied",
			fen:    "8/8/8/8/2k5/8/8/4K1NR w K - 0 1",
			kingSq: E1,
			rookSq: H1,
			toFile: FileG,
			dir:    1,
			want:   false,
		},
		{
			name:   "king in check",
			fen:    "4k3/8/8/8/8/4r3/8/4K2R w K - 0 1",
			kingSq: E1,
			rookSq: H1,
			toFile: FileG,
			dir:    1,
			want:   false,
		},
		{
			name:   "king passes through attacked square",
			fen:    "4k3/8/8/8/8/5r2/8/4K2R w K - 0 1",
			kingSq: E1,
			rookSq: H1,
			toFile: FileG,
			dir:    1,
			want:   false,
		},
		{
			name:   "kinside destination square attacked",
			fen:    "4k3/8/8/8/8/6r1/8/4K2R w K - 0 1",
			kingSq: E1,
			rookSq: H1,
			toFile: FileG,
			dir:    1,
			want:   false,
		},
		{
			name:   "clear queenside path",
			fen:    "4k3/8/8/8/8/8/8/R3K3 w Q - 0 1",
			kingSq: E1,
			rookSq: A1,
			toFile: FileC,
			dir:    -1,
			want:   true,
		},
		{
			name:   "blocked queenside path",
			fen:    "3k4/8/8/8/8/8/8/RN2K3 w Q - 0 2",
			kingSq: E1,
			rookSq: A1,
			toFile: FileC,
			dir:    -1,
			want:   false,
		},
		{
			name:   "queenside destination square attacked",
			fen:    "4k3/8/8/8/8/2r5/8/R3K3 w Q - 0 1",
			kingSq: E1,
			rookSq: A1,
			toFile: FileC,
			dir:    -1,
			want:   false,
		},
		{
			name:   "clear black kingside path",
			fen:    "4k2r/8/8/8/8/8/8/4K3 b k - 0 1",
			kingSq: E8,
			rookSq: H8,
			toFile: FileG,
			dir:    1,
			want:   true,
		},
		{
			name:   "black king passes through attacked square",
			fen:    "r3k3/8/8/3R4/8/8/8/4K3 w q - 0 1",
			kingSq: E8,
			rookSq: A8,
			toFile: FileC,
			dir:    -1,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			got := canCastle(&pos, tt.kingSq, tt.rookSq)
			if got != tt.want {
				t.Errorf("expected %v but got %v\n%v", got, tt.want, pos.String())
			}
		})
	}
}

func TestCanCastle_Chess960(t *testing.T) {
	withChess960(t)
	tests := []struct {
		name   string
		fen    string
		kingSq Square
		rookSq Square
		want   bool
	}{
		{
			name:   "clear kingside path",
			fen:    "8/8/8/8/2k5/8/8/1K5R w H - 0 1",
			kingSq: B1,
			rookSq: H1,
			want:   true,
		},
		{
			name:   "castling path blocked",
			fen:    "8/8/8/8/2k5/8/8/1K1N3R w H - 0 1",
			kingSq: B1,
			rookSq: H1,
			want:   false,
		},
		{
			name:   "destination square occupied",
			fen:    "8/8/8/3b4/2k5/8/8/1K3RN1 w F - 0 1",
			kingSq: B1,
			rookSq: F1,
			want:   false,
		},
		{
			name:   "king passes through attacked square",
			fen:    "8/8/8/8/2k5/8/1n6/1K5R w H - 0 1",
			kingSq: B1,
			rookSq: H1,
			want:   false,
		},
		{
			name:   "kingside destination square attacked",
			fen:    "8/8/8/8/2k5/5n2/8/1K5R w H - 0 1",
			kingSq: B1,
			rookSq: H1,
			want:   false,
		},
		{
			name:   "king jumps over rook kingside",
			fen:    "8/8/8/8/1k6/8/8/1KR5 w C - 0 1",
			kingSq: B1,
			rookSq: C1,
			want:   true,
		},
		{
			name:   "king jumps over rook queenside",
			fen:    "8/8/8/8/1k6/8/8/5RK1 w F - 0 1",
			kingSq: G1,
			rookSq: F1,
			want:   true,
		},
		{
			name:   "king jumps over rook queenside, but destination square attacked",
			fen:    "8/8/8/6b1/1k6/8/8/5RK1 w F - 0 1",
			kingSq: G1,
			rookSq: F1,
			want:   false,
		},
		{
			name:   "king and rook move right, but path attacked",
			fen:    "8/8/8/8/1k6/2b5/8/1KR5 w C - 0 1",
			kingSq: B1,
			rookSq: C1,
			want:   false,
		},
		{
			name:   "rook jumps over king kingside",
			fen:    "8/8/8/8/1k6/8/8/6KR w H - 0 1",
			kingSq: G1,
			rookSq: H1,
			want:   true,
		},
		{
			name:   "rook jumps over king queenside",
			fen:    "8/8/8/8/1k1b4/8/8/RK6 w A - 0 1",
			kingSq: B1,
			rookSq: A1,
			want:   true,
		},
		{
			name:   "rook jumps over king but rook destination square blocked",
			fen:    "8/8/8/8/1k6/8/8/5NKR w H - 0 1",
			kingSq: G1,
			rookSq: H1,
			want:   false,
		},
		{
			name:   "rook and king swap positions",
			fen:    "8/8/8/8/1k6/8/8/5KR1 w G - 0 1",
			kingSq: F1,
			rookSq: G1,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad FEN: %v", err)
			}

			got := canCastle(&pos, tt.kingSq, tt.rookSq)
			if got != tt.want {
				t.Errorf("expected %v but got %v\n%v", tt.want, got, pos.String())
			}
		})
	}
}

func BenchmarkGenMoves(b *testing.B) {
	fens := []string{
		StartingFEN,
		"r1bqkbnr/pppp1ppp/2n5/4p3/1P6/2N5/P1PPPPPP/R1BQKBNR w KQkq - 0 1",
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
	}

	positions := make([]*Position, 0, len(fens))
	for _, fen := range fens {
		pos, err := ParseFEN(fen)
		if err != nil {
			b.Fatalf("Bad FEN: %v", fen)
		}
		positions = append(positions, &pos)
	}

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		var movelist Movelist
		pos := positions[i%len(positions)]
		GenPawnMoves(pos, &movelist)
	}
}
