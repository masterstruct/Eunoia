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
			fen:     StartingFEN,
			move:    NewMove(E2, E3),
			wantFEN: "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			name:    "white quiet move clears en passant",
			fen:     "rnbqkb1r/pp1ppppp/5n2/2pP4/8/8/PPP1PPPP/RNBQKBNR w KQkq c6 0 3",
			move:    NewMove(G1, F3),
			wantFEN: "rnbqkb1r/pp1ppppp/5n2/2pP4/8/5N2/PPP1PPPP/RNBQKB1R b KQkq - 1 3",
		},
		{
			name:    "white double push does not set en passant",
			fen:     StartingFEN,
			move:    NewDoublePush(E2, E4),
			wantFEN: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
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
			fen:     "k7/8/8/3p4/4P3/8/8/K7 b - - 0 1",
			move:    NewCapture(D5, E4),
			wantFEN: "k7/8/8/8/4p3/8/8/K7 w - - 0 2",
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
			name:    "white kingside castle moves rook too",
			fen:     "4k3/8/8/8/8/8/8/4K2R w K - 0 1",
			move:    NewCastle(E1, H1),
			wantFEN: "4k3/8/8/8/8/8/8/5RK1 b - - 1 1",
		},
		{
			name:    "white queenside castle moves rook too",
			fen:     "4k3/8/8/8/8/8/8/R3K3 w Q - 0 1",
			move:    NewCastle(E1, A1),
			wantFEN: "4k3/8/8/8/8/8/8/2KR4 b - - 1 1",
		},
		{
			name:    "black kingside castle moves rook too",
			fen:     "4k2r/8/8/8/8/8/8/4K3 b k - 1 1",
			move:    NewCastle(E8, H8),
			wantFEN: "5rk1/8/8/8/8/8/8/4K3 w - - 2 2",
		},
		{
			name:    "black queenside castle moves rook too",
			fen:     "r3k3/8/8/8/8/8/8/4K3 b q - 16 160",
			move:    NewCastle(E8, A8),
			wantFEN: "2kr4/8/8/8/8/8/8/4K3 w - - 17 161",
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
			name:    "capturing black rook on its home square forfeits black's castling right",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:    NewCapture(A1, A8),
			wantFEN: "R3k2r/8/8/8/8/8/8/4K2R b Kk - 0 1",
		},
		{
			name:    "capturing white rook on its home square forfeits white's castling right",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			move:    NewCapture(H8, H1),
			wantFEN: "r3k3/8/8/8/8/8/8/R3K2r w Qq - 0 2",
		},
		{
			name:    "moving a rook from its home square forfeits castling right",
			fen:     "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			move:    NewMove(A1, D1),
			wantFEN: "r3k2r/8/8/8/8/8/8/3RK2R b Kkq - 1 1",
		},
		{
			name:    "white promotion places knight not pawn",
			fen:     "8/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewPromo(E7, E8, Knight),
			wantFEN: "4N3/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "white promotion places bishop not pawn",
			fen:     "8/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewPromo(E7, E8, Bishop),
			wantFEN: "4B3/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "white promotion places rook not pawn",
			fen:     "8/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewPromo(E7, E8, Rook),
			wantFEN: "4R3/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "white promotion places queen not pawn",
			fen:     "8/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewPromo(E7, E8, Queen),
			wantFEN: "4Q3/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "black promotion places knight not pawn",
			fen:     "8/2K5/8/6k1/8/8/2p5/8 b - - 0 1",
			move:    NewPromo(C2, C1, Knight),
			wantFEN: "8/2K5/8/6k1/8/8/8/2n5 w - - 0 2",
		},
		{
			name:    "black promotion places bishop not pawn",
			fen:     "8/2K5/8/6k1/8/8/2p5/8 b - - 0 1",
			move:    NewPromo(C2, C1, Bishop),
			wantFEN: "8/2K5/8/6k1/8/8/8/2b5 w - - 0 2",
		},
		{
			name:    "black promotion places rook not pawn",
			fen:     "8/2K5/8/6k1/8/8/2p5/8 b - - 0 1",
			move:    NewPromo(C2, C1, Rook),
			wantFEN: "8/2K5/8/6k1/8/8/8/2r5 w - - 0 2",
		},
		{
			name:    "black promotion places queen not pawn",
			fen:     "8/2K5/8/6k1/8/8/2p5/8 b - - 0 1",
			move:    NewPromo(C2, C1, Queen),
			wantFEN: "8/2K5/8/6k1/8/8/8/2q5 w - - 0 2",
		},
		{
			name:    "white capture promotion removes victim and places queen",
			fen:     "3r4/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewCapturePromo(E7, D8, Queen),
			wantFEN: "3Q4/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "white capture promotion removes victim and places knight",
			fen:     "3r4/4P3/8/8/8/8/8/4K2k w - - 0 1",
			move:    NewCapturePromo(E7, D8, Knight),
			wantFEN: "3N4/8/8/8/8/8/8/4K2k b - - 0 1",
		},
		{
			name:    "black capture promotion removes victim, places bishop",
			fen:     "4K2k/8/8/8/8/8/4p3/3N4 b - - 0 1",
			move:    NewCapturePromo(E2, D1, Bishop),
			wantFEN: "4K2k/8/8/8/8/8/8/3b4 w - - 0 2",
		},
		{
			name:    "black capture promotion removes victim, places rook",
			fen:     "4K2k/8/8/8/8/8/4p3/3N4 b - - 0 1",
			move:    NewCapturePromo(E2, D1, Rook),
			wantFEN: "4K2k/8/8/8/8/8/8/3r4 w - - 0 2",
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
			move:    NewCastle(E8, H8),
			wantFEN: "5rk1/8/8/8/8/8/8/4K3 w - - 1 2",
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
			name:    "black captures white rook on h1, clears K",
			fen:     "4k3/8/8/8/7r/8/8/R3K2R b KQkq - 0 1",
			move:    NewCapture(H4, H1),
			wantFEN: "4k3/8/8/8/8/8/8/R3K2r w Qkq - 0 2",
		},
		{
			name:    "rook move from non-home square doesn't touch castling rights",
			fen:     "4k3/8/8/8/8/8/3R4/R3K2R w KQ - 0 1",
			move:    NewMove(D2, D5),
			wantFEN: "4k3/8/8/3R4/8/8/8/R3K2R b KQ - 1 1",
		},
		{
			name:    "capturing a rook off its home square doesn't touch castling rights",
			fen:     "4k3/8/8/3r4/1N6/8/8/R3K2R w KQ - 0 1",
			move:    NewCapture(B4, D5),
			wantFEN: "4k3/8/8/3N4/8/8/8/R3K2R b KQ - 0 1",
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
}

func BenchmarkMakeMove(b *testing.B) {
	tests := []struct {
		fen      string
		move     Move
		chess960 bool
	}{
		{StartingFEN, NewMove(E2, E3), false},
		{"rnbqkb1r/pp1ppppp/5n2/2pP4/8/8/PPP1PPPP/RNBQKBNR w KQkq c6 0 3", NewMove(G1, F3), false},
		{StartingFEN, NewDoublePush(E2, E4), false},
		{"rnbqkb1r/pppppppp/5n2/3P4/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 2", NewDoublePush(C7, C5), false},
		{"8/2p5/3p4/KP5r/1R3p1k/8/6P1/8 w - - 0 1", NewDoublePush(G2, G4), false},
		{"k7/8/8/3p4/4P3/8/8/K7 b - - 0 1", NewCapture(D5, E4), false},
		{"k7/8/8/3p4/4P3/8/8/K7 w - - 17 5", NewCapture(E4, D5), false},
		{"4k3/8/8/8/8/8/8/R3K3 w - - 3 5", NewMove(A1, A5), false},
		{"4k3/8/8/8/8/8/4P3/4K3 w - - 9 5", NewMove(E2, E3), false},
		{"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1", NewEnPassant(E5, D6), false},
		{"4k3/8/8/8/8/8/8/4K2R w K - 0 1", NewCastle(E1, H1), false},
		{"4k3/8/8/8/8/8/8/R3K3 w Q - 0 1", NewCastle(E1, A1), false},
		{"4k2r/8/8/8/8/8/8/4K3 b k - 1 1", NewCastle(E8, H8), false},
		{"r3k3/8/8/8/8/8/8/4K3 b q - 16 160", NewCastle(E8, A8), false},
		{"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1", NewMove(E1, E2), false},
		{"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1", NewMove(H1, H2), false},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", NewCapture(A1, A8), false},
		{"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1", NewCapture(H8, H1), false},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", NewMove(A1, D1), false},
		{"8/4P3/8/8/8/8/8/4K2k w - - 0 1", NewPromo(E7, E8, Knight), false},
		{"8/4P3/8/8/8/8/8/4K2k w - - 0 1", NewPromo(E7, E8, Bishop), false},
		{"8/4P3/8/8/8/8/8/4K2k w - - 0 1", NewPromo(E7, E8, Rook), false},
		{"8/4P3/8/8/8/8/8/4K2k w - - 0 1", NewPromo(E7, E8, Queen), false},
		{"8/2K5/8/6k1/8/8/2p5/8 b - - 0 1", NewPromo(C2, C1, Knight), false},
		{"8/2K5/8/6k1/8/8/2p5/8 b - - 0 1", NewPromo(C2, C1, Bishop), false},
		{"8/2K5/8/6k1/8/8/2p5/8 b - - 0 1", NewPromo(C2, C1, Rook), false},
		{"8/2K5/8/6k1/8/8/2p5/8 b - - 0 1", NewPromo(C2, C1, Queen), false},
		{"3r4/4P3/8/8/8/8/8/4K2k w - - 0 1", NewCapturePromo(E7, D8, Queen), false},
		{"3r4/4P3/8/8/8/8/8/4K2k w - - 0 1", NewCapturePromo(E7, D8, Knight), false},
		{"4K2k/8/8/8/8/8/4p3/3N4 b - - 0 1", NewCapturePromo(E2, D1, Bishop), false},
		{"4K2k/8/8/8/8/8/4p3/3N4 b - - 0 1", NewCapturePromo(E2, D1, Rook), false},
		{"8/8/8/8/8/8/8/4K2k b - - 0 7", NewMove(H1, H2), false},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", NewMove(G8, F6), false},
		{"4k3/8/8/8/3pP3/8/8/4K3 b - e3 0 1", NewEnPassant(D4, E3), false},
		{"4k3/8/8/2Pp4/8/8/8/4K3 w - d6 0 1", NewMove(E1, E2), false},
		{"4k2r/8/8/8/8/8/8/4K3 b k - 0 1", NewCastle(E8, H8), false},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", NewMove(E1, E2), false},
		{"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1", NewMove(E8, E7), false},
		{"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1", NewMove(H8, H7), false},
		{"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1", NewMove(A8, A6), false},
		{"4k3/8/8/8/7r/8/8/R3K2R b KQkq - 0 1", NewCapture(H4, H1), false},
		{"4k3/8/8/8/8/8/8/R2RK2R w KQ - 0 1", NewMove(D1, D5), false},
		{"4k3/8/8/3r4/1N6/8/8/R3K2R w KQ - 0 1", NewCapture(B4, D5), false},

		{"qb3rk1/pp3rp1/2p1p2p/2P4P/3PQ1P1/1PN5/P4P2/3R1K1R w DH - 1 19", NewCastle(F1, H1), true},
		{"bqrbkrR1/4pp1Q/8/1pp2N2/8/2PnP3/1P5P/B1RBKN2 w Ccf - 2 15", NewMove(E1, D2), true},
		{"2rb1R2/2k1ppN1/2b5/8/1Pp5/2P1P3/7P/B1KB1N2 b - b3 0 22", NewEnPassant(C4, B3), true},
		{"rkr1q1b1/pp1p1ppp/2p3n1/4p3/4P3/2P1NPN1/PP1P2PP/R1KB1Q2 b ac - 0 10", NewMove(C8, C7), true},
		{"rk2q1b1/pprp1ppp/2p3n1/4p3/4P3/2P1NPN1/PP1P2PP/R1KB1Q2 b a - 2 12", NewCastle(B8, A8), true},
		{"4k3/8/8/8/8/8/8/5KRQ w G - 0 1", NewCastle(F1, G1), true},
		{"2r3kr/8/8/8/8/8/8/RK1R4 b hc - 1 1", NewCastle(G8, H8), true},
		{"3r2kr/8/8/8/8/8/8/RK3R2 w AF - 0 1", NewCastle(B1, A1), true},
		{"1r1kr3/8/8/8/8/4R3/2K1R3/8 w bf - 0 1", NewCapture(E3, E8), true},
	}

	type fixture struct {
		pos      Position
		move     Move
		chess960 bool
	}

	fixtures := make([]fixture, len(tests))
	for i, tt := range tests {
		pos, err := ParseFEN(tt.fen)
		if err != nil {
			b.Fatalf("bad benchmark FEN: %v", err)
		}
		fixtures[i] = fixture{pos: pos, move: tt.move, chess960: tt.chess960}
	}

	for i := 0; b.Loop(); i++ {
		f := fixtures[i%len(fixtures)]
		if f.chess960 {
			SetChess960(true)
		}
		_ = f.pos.MakeMove(f.move)
		if f.chess960 {
			SetChess960(false)
		}
	}
}

func TestMakeMove_Chess960(t *testing.T) {
	withChess960(t)
	tests := []struct {
		name    string
		fen     string
		move    Move
		wantFEN string // resulting position, roundtrip through FEN for readability
	}{
		{
			name:    "White castle, king already on destination square",
			fen:     "qb3rk1/pp3rp1/2p1p2p/2P4P/3PQ1P1/1PN5/P4P2/3R1K1R w DH - 1 19",
			move:    NewCastle(F1, H1),
			wantFEN: "qb3rk1/pp3rp1/2p1p2p/2P4P/3PQ1P1/1PN5/P4P2/3R1RK1 b - - 2 19",
		},
		{
			name:    "White king moves, clear castling rights",
			fen:     "bqrbkrR1/4pp1Q/8/1pp2N2/8/2PnP3/1P5P/B1RBKN2 w Ccf - 2 15",
			move:    NewMove(E1, D2),
			wantFEN: "bqrbkrR1/4pp1Q/8/1pp2N2/8/2PnP3/1P1K3P/B1RB1N2 b cf - 3 15",
		},
		{
			name:    "En passant",
			fen:     "2rb1R2/2k1ppN1/2b5/8/1Pp5/2P1P3/7P/B1KB1N2 b - b3 0 22",
			move:    NewEnPassant(C4, B3),
			wantFEN: "2rb1R2/2k1ppN1/2b5/8/8/1pP1P3/7P/B1KB1N2 w - - 0 23",
		},
		{
			name:    "Black rook moves, clear castling right",
			fen:     "rkr1q1b1/pp1p1ppp/2p3n1/4p3/4P3/2P1NPN1/PP1P2PP/R1KB1Q2 b ac - 0 10",
			move:    NewMove(C8, C7),
			wantFEN: "rk2q1b1/pprp1ppp/2p3n1/4p3/4P3/2P1NPN1/PP1P2PP/R1KB1Q2 w a - 1 11",
		},
		{
			name:    "Black castle",
			fen:     "rk2q1b1/pprp1ppp/2p3n1/4p3/4P3/2P1NPN1/PP1P2PP/R1KB1Q2 b a - 2 12",
			move:    NewCastle(B8, A8),
			wantFEN: "2krq1b1/pprp1ppp/2p3n1/4p3/4P3/2P1NPN1/PP1P2PP/R1KB1Q2 w - - 3 13",
		},
		{
			name:    "White castle, king and rook swap, queen in corner",
			fen:     "4k3/8/8/8/8/8/8/5KRQ w G - 0 1",
			move:    NewCastle(F1, G1),
			wantFEN: "4k3/8/8/8/8/8/8/5RKQ b - - 1 1",
		},
		{
			name:    "Black castle, edge of board",
			fen:     "2r3kr/8/8/8/8/8/8/RK1R4 b hc - 1 1",
			move:    NewCastle(G8, H8),
			wantFEN: "2r2rk1/8/8/8/8/8/8/RK1R4 w - - 2 2",
		},
		{
			name:    "White castle, edge of board",
			fen:     "3r2kr/8/8/8/8/8/8/RK3R2 w AF - 0 1",
			move:    NewCastle(B1, A1),
			wantFEN: "3r2kr/8/8/8/8/8/8/2KR1R2 b - - 1 1",
		},
		{
			name:    "Capturing black rook clears correct castling right",
			fen:     "1r1kr3/8/8/8/8/4R3/2K1R3/8 w bf - 0 1",
			move:    NewCapture(E3, E8),
			wantFEN: "1r1kR3/8/8/8/8/8/2K1R3/8 b b - 0 1",
		},
		{
			name:    "Capturing black rook clears correct white AND black castling right",
			fen:     "1r1k2r1/8/8/8/8/8/8/1R1K2R1 w BGbg - 1 1",
			move:    NewCapture(B1, B8),
			wantFEN: "1R1k2r1/8/8/8/8/8/8/3K2R1 b Gg - 0 1",
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
