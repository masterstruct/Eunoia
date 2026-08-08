package board

import "testing"

func TestMakeMove(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		move    Move
		wantFEN string // resulting position, roundtrip through FEN for readability
	}{
		{
			name:    "quiet pawn push clears en passant",
			fen:     startingFEN,
			move:    NewMove(E2, E3),
			wantFEN: "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name:    "double push sets en passant unconditionally",
			fen:     startingFEN,
			move:    NewDoublePush(E2, E4),
			wantFEN: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			name:    "en passant sets ep square even though pinned pawn can't capture",
			fen:     "8/2p5/3p4/KP5r/1R3p1k/8/6P1/8 w - - 0 1",
			move:    NewDoublePush(G2, G4),
			wantFEN: "8/2p5/3p4/KP5r/1R3pPk/8/8/8 b - g3 0 1",
		},
		{
			name:    "capture removes victim",
			fen:     "k7/8/8/3p4/4P3/8/8/K7 w - - 0 1",
			move:    NewCapture(E4, D5),
			wantFEN: "k7/8/8/3P4/8/8/8/K7 b - - 0 1",
		},
		{
			name:    "capture resets halfmove clock",
			fen:     "k7/8/8/3p4/4P3/8/8/K7 w - - 17 5",
			move:    NewCapture(E4, D5),
			wantFEN: "k7/8/8/3P4/8/8/8/K7 b - - 0 5",
		},
		{
			name:    "quiet non-pawn move increments halfmove clock",
			fen:     "4k3/8/8/8/8/8/8/R3K3 w - - 3 5",
			move:    NewMove(A1, A5),
			wantFEN: "4k3/8/8/R7/8/8/8/4K3 b - - 4 5",
		},
		{
			name:    "pawn push resets halfmove clock",
			fen:     "4k3/8/8/8/8/8/4P3/4K3 w - - 9 5",
			move:    NewMove(E2, E3),
			wantFEN: "4k3/8/8/8/8/4P3/8/4K3 b - - 0 5",
		},
		{
			name:    "en passant capture removes correct pawn, not destination square",
			fen:     "4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
			move:    NewEnPassant(E5, D6),
			wantFEN: "4k3/8/3P4/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:    "kingside castle moves rook too",
			fen:     "4k3/8/8/8/8/8/8/4K2R w K - 0 1",
			move:    NewCastle(E1, G1, true),
			wantFEN: "4k3/8/8/8/8/8/8/5RK1 b - - 1 1",
		},
		{
			name:    "queenside castle moves rook too",
			fen:     "4k3/8/8/8/8/8/8/R3K3 w Q - 0 1",
			move:    NewCastle(E1, C1, false),
			wantFEN: "4k3/8/8/8/8/8/8/2KR4 b - - 1 1",
		},
		{
			name:    "king move forfeits both castling rights",
			fen:     "4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			move:    NewMove(E1, E2),
			wantFEN: "4k3/8/8/8/8/8/4K3/R6R b - - 1 1",
		},
		{
			name:    "kingside rook move forfeits only kingside right",
			fen:     "4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			move:    NewMove(H1, H2),
			wantFEN: "4k3/8/8/8/8/8/7R/R3K3 b Q - 1 1",
		},
		{
			name:    "capturing enemy rook on its home square forfeits that side's right",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:    NewCapture(A1, A8),
			wantFEN: "R3k2r/8/8/8/8/8/8/4K2R b Kk - 0 1",
		},
		{
			name:    "promotion places promoted piece not pawn",
			fen:     "8/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewPromo(E7, E8, Queen),
			wantFEN: "4Q3/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "capture promotion removes victim and places promoted piece",
			fen:     "3r4/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewCapturePromo(E7, D8, Queen),
			wantFEN: "3Q4/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "ply always increments regardless of move type",
			fen:     "8/8/8/8/8/8/8/4K2k b - - 0 7",
			move:    NewMove(H1, H2),
			wantFEN: "8/8/8/8/8/8/7k/4K3 w - - 1 8",
		},
		// TODO: add more tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad test FEN: %v", err)
			}

			got := pos.MakeMove(tt.move)

			want, err := ParseFEN(tt.wantFEN)
			if err != nil {
				t.Fatalf("failed to generate FEN: %v", err)
			}

			if got != want {
				t.Errorf("MakeMove(%v) on %q\ngot %v:\n%v\nwant: %v\n%v", tt.move, tt.fen, got.FEN(), got, want.FEN(), want)
			}
		})
	}
}
