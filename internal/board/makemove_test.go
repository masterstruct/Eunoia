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
			name:    "quiet move",
			fen:     startingFEN,
			move:    NewMove(E2, E3),
			wantFEN: "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name:    "quiet move clears en passant",
			fen:     "rnbqkb1r/pp1ppppp/5n2/2pP4/8/8/PPP1PPPP/RNBQKBNR w KQkq c6 0 3",
			move:    NewMove(G1, F3),
			wantFEN: "rnbqkb1r/pp1ppppp/5n2/2pP4/8/5N2/PPP1PPPP/RNBQKB1R b KQkq - 1 3",
		},
		{
			name:    "white double push sets en passant",
			fen:     startingFEN,
			move:    NewDoublePush(E2, E4),
			wantFEN: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		{
			name:    "black double push sets en passant",
			fen:     "rnbqkb1r/pppppppp/5n2/3P4/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 2",
			move:    NewDoublePush(C7, C5),
			wantFEN: "rnbqkb1r/pp1ppppp/5n2/2pP4/8/8/PPP1PPPP/RNBQKBNR w KQkq c6 0 3",
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
			name:    "moving a rook from its home square forfeits castling right",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:    NewMove(A1, A3),
			wantFEN: "r3k2r/8/8/8/8/R7/8/4K2R b Kkq - 1 1",
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
		{
			name:    "black quiet move clears en passant",
			fen:     "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			move:    NewMove(G8, F6),
			wantFEN: "rnbqkb1r/pppppppp/5n2/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 1 2",
		},
		{
			name:    "black captures en passant, removes correct pawn",
			fen:     "4k3/8/8/8/3pP3/8/8/4K3 b - e3 0 1",
			move:    NewEnPassant(D4, E3),
			wantFEN: "4k3/8/8/8/8/4p3/8/4K3 w - - 0 2",
		},
		{
			name:    "any move clears stale en passant, not just the capturing pawn",
			fen:     "4k3/8/8/2Pp4/8/8/8/4K3 w - d6 0 1",
			move:    NewMove(E1, E2),
			wantFEN: "4k3/8/8/2Pp4/8/8/4K3/8 b - - 1 1",
		},
		{
			name:    "black kingside castle moves rook too",
			fen:     "4k2r/8/8/8/8/8/8/4K3 b k - 0 1",
			move:    NewCastle(E8, G8, true),
			wantFEN: "5rk1/8/8/8/8/8/8/4K3 w - - 1 2",
		},
		{
			name:    "black queenside castle moves rook too",
			fen:     "r3k3/8/8/8/8/8/8/4K3 b q - 0 1",
			move:    NewCastle(E8, C8, false),
			wantFEN: "2kr4/8/8/8/8/8/8/4K3 w - - 1 2",
		},
		{
			name:    "white king move clears KQ only, kq untouched",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:    NewMove(E1, E2),
			wantFEN: "r3k2r/8/8/8/8/8/4K3/R6R b kq - 1 1",
		},
		{
			name:    "black king move clears kq only, KQ untouched",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			move:    NewMove(E8, E7),
			wantFEN: "r6r/4k3/8/8/8/8/8/R3K2R w KQ - 1 2",
		},
		{
			name:    "black kingside rook move clears k only, Qq untouched",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			move:    NewMove(H8, H7),
			wantFEN: "r3k3/7r/8/8/8/8/8/R3K2R w KQq - 1 2",
		},
		{
			name:    "black queenside rook move clears q only, Kk untouched",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			move:    NewMove(A8, A6),
			wantFEN: "4k2r/8/r7/8/8/8/8/R3K2R w KQk - 1 2",
		},
		{
			name:    "black captures white rook on h1, clears K via captured-square not mover-square",
			fen:     "4k3/8/8/8/7r/8/8/R3K2R b KQkq - 0 1",
			move:    NewCapture(H4, H1),
			wantFEN: "4k3/8/8/8/8/8/8/R3K2r w Qkq - 0 2",
		},
		{
			name:    "rook move from non-home square doesn't touch castling rights",
			fen:     "4k3/8/8/8/8/8/8/R2RK2R w KQ - 0 1",
			move:    NewMove(D1, D5),
			wantFEN: "4k3/8/8/3R4/8/8/8/R3K2R b KQ - 1 1",
		},
		{
			name:    "capturing a rook off its home square doesn't touch castling rights",
			fen:     "4k3/8/8/3r4/1N6/8/8/R3K2R w KQ - 0 1",
			move:    NewCapture(B4, D5),
			wantFEN: "4k3/8/8/3N4/8/8/8/R3K2R b KQ - 0 1",
		},
		{
			name:    "underpromotion places knight, not queen",
			fen:     "8/4P3/8/8/8/8/8/4K2k w - - 12 5",
			move:    NewPromo(E7, E8, Knight),
			wantFEN: "4N3/8/8/8/8/8/8/4K2k b - - 0 5",
		},
		{
			name:    "black promotion places correct piece",
			fen:     "4K2k/8/8/8/8/8/4p3/8 b - - 0 1",
			move:    NewPromo(E2, E1, Rook),
			wantFEN: "4K2k/8/8/8/8/8/8/4r3 w - - 0 2",
		},
		{
			name:    "black capture promotion removes victim, places promoted piece",
			fen:     "4K2k/8/8/8/8/8/4p3/3N4 b - - 0 1",
			move:    NewCapturePromo(E2, D1, Queen),
			wantFEN: "4K2k/8/8/8/8/8/8/3q4 w - - 0 2",
		},
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

			makemove_assertPositionEqual(t, got, want)
		})
	}
}

func makemove_assertPositionEqual(t *testing.T, got, want Position) {
	t.Helper()

	if got.FEN() != want.FEN() {
		t.Errorf("FEN mismatch\ngot: %s\n%v\nwant: %s\n%v", got.FEN(), got, want.FEN(), want)
	}
	if got.SideToMove != want.SideToMove {
		t.Errorf("SideToMove mismatch\ngot: %v\nwant: %v", got.SideToMove, want.SideToMove)
	}
	if got.CastlingRights != want.CastlingRights {
		t.Errorf("CastlingRights mismatch\ngot: %v\nwant: %v", got.CastlingRights, want.CastlingRights)
	}
	if got.EnPassant != want.EnPassant {
		t.Errorf("EnPassant mismatch\ngot: %v\nwant: %v", got.EnPassant, want.EnPassant)
	}
	if got.HalfmoveClock != want.HalfmoveClock {
		t.Errorf("HalfmoveClock mismatch\ngot: %v\nwant: %v", got.HalfmoveClock, want.HalfmoveClock)
	}
	if got.Ply != want.Ply {
		t.Errorf("Ply mismatch\ngot: %v\nwant: %v", got.Ply, want.Ply)
	}
	if got.Hash != want.Hash {
		t.Errorf("Hash mismatch\ngot: %v\nwant: %v", got.Hash, want.Hash)
	}
}
