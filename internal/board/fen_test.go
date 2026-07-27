package board

import (
	"errors"
	"testing"
)

func wantPosition(placements []placement, sideToMove Color, castling CastlingRights, epSq Square, halfmove uint8, ply uint16) Position {
	pos := NewPosition()

	pos.SideToMove = sideToMove
	pos.CastlingRights = castling
	pos.EnPassant = epSq
	pos.HalfmoveClock = halfmove
	pos.Ply = ply

	for _, pl := range placements {
		pos.PlacePiece(pl.p, pl.sq)
	}
	return pos
}

func assertPositionEqual(t *testing.T, got, want Position, fen string) {
	t.Helper()

	if got.Board != want.Board {
		t.Errorf("%q: board mismatch\ngot:\n%v\nwant: \n%v", fen, got, want)
	}
	if got.pieces != want.pieces {
		t.Errorf("%q: piece bitboards mismatch\ngot:  %v\nwant: %v", fen, got.pieces, want.pieces)
	}
	if got.colors != want.colors {
		t.Errorf("%q: color bitboards mismatch\ngot:  %v\nwant: %v", fen, got.colors, want.colors)
	}
	if got.SideToMove != want.SideToMove {
		t.Errorf("%q: side to move: got %v want %v", fen, got.SideToMove, want.SideToMove)
	}
	if got.CastlingRights != want.CastlingRights {
		t.Errorf("%q: castling rights: got %v want %v", fen, got.CastlingRights, want.CastlingRights)
	}
	if got.EnPassant != want.EnPassant {
		t.Errorf("%q: en passant: got %v want %v", fen, got.EnPassant, want.EnPassant)
	}
	if got.HalfmoveClock != want.HalfmoveClock {
		t.Errorf("%q: halfmove clock: got %v want %v", fen, got.HalfmoveClock, want.HalfmoveClock)
	}
	if got.Ply != want.Ply {
		t.Errorf("%q: ply: got %v want %v", fen, got.Ply, want.Ply)
	}
}

func TestParseFEN_Valid(t *testing.T) {
	tests := []struct {
		name       string
		fen        string
		placements []placement
		sideToMove Color
		castling   CastlingRights
		epSq       Square
		halfmove   uint8
		ply        uint16
	}{
		{
			name: "starting position",
			fen:  startingFEN,
			placements: []placement{
				{BlackRook, A8}, {BlackKnight, B8}, {BlackBishop, C8}, {BlackQueen, D8},
				{BlackKing, E8}, {BlackBishop, F8}, {BlackKnight, G8}, {BlackRook, H8},

				{BlackPawn, A7}, {BlackPawn, B7}, {BlackPawn, C7}, {BlackPawn, D7},
				{BlackPawn, E7}, {BlackPawn, F7}, {BlackPawn, G7}, {BlackPawn, H7},

				{WhitePawn, A2}, {WhitePawn, B2}, {WhitePawn, C2}, {WhitePawn, D2},
				{WhitePawn, E2}, {WhitePawn, F2}, {WhitePawn, G2}, {WhitePawn, H2},

				{WhiteRook, A1}, {WhiteKnight, B1}, {WhiteBishop, C1}, {WhiteQueen, D1},
				{WhiteKing, E1}, {WhiteBishop, F1}, {WhiteKnight, G1}, {WhiteRook, H1},
			},
			sideToMove: White, castling: AllCastling, epSq: NoSquare, halfmove: 0, ply: 0,
		},
		{
			name:       "empty board",
			fen:        "8/8/8/8/8/8/8/8 w - - 0 1",
			placements: nil,
			sideToMove: White, castling: NoCastling, epSq: NoSquare, halfmove: 0, ply: 0,
		},
		{
			name: "en passant set",
			fen:  "rnbqkb1r/pppp1ppp/5n2/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
			placements: []placement{
				{BlackRook, A8}, {BlackKnight, B8}, {BlackBishop, C8}, {BlackQueen, D8},
				{BlackKing, E8}, {BlackBishop, F8}, {BlackRook, H8},

				{BlackPawn, A7}, {BlackPawn, B7}, {BlackPawn, C7}, {BlackPawn, D7},
				{BlackPawn, F7}, {BlackPawn, G7}, {BlackPawn, H7},

				{BlackKnight, F6},

				{WhitePawn, D5}, {BlackPawn, E5},

				{WhitePawn, A2}, {WhitePawn, B2}, {WhitePawn, C2}, {WhitePawn, E2},
				{WhitePawn, F2}, {WhitePawn, G2}, {WhitePawn, H2}, {WhitePawn, D5},

				{WhiteRook, A1}, {WhiteKnight, B1}, {WhiteBishop, C1}, {WhiteQueen, D1},
				{WhiteKing, E1}, {WhiteBishop, F1}, {WhiteKnight, G1}, {WhiteRook, H1},
			},
			sideToMove: White, castling: AllCastling, epSq: E6, halfmove: 0, ply: 4,
		},

		{
			name:       "single piece each corner",
			fen:        "K6k/8/8/8/8/8/8/n6B w - - 0 1",
			placements: []placement{{WhiteKing, A8}, {BlackKing, H8}, {BlackKnight, A1}, {WhiteBishop, H1}},
			sideToMove: White, castling: NoCastling, epSq: NoSquare, halfmove: 0, ply: 0,
		},
		{
			name: "kiwipete",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -",
			placements: []placement{
				{BlackRook, A8}, {BlackKing, E8}, {BlackRook, H8},
				{BlackPawn, A7}, {BlackPawn, C7}, {BlackPawn, D7}, {BlackQueen, E7}, {BlackPawn, F7}, {BlackBishop, G7},
				{BlackBishop, A6}, {BlackKnight, B6}, {BlackPawn, E6}, {BlackKnight, F6}, {BlackPawn, G6},
				{WhitePawn, D5}, {WhiteKnight, E5},
				{BlackPawn, B4}, {WhitePawn, E4},
				{WhiteKnight, C3}, {WhiteQueen, F3}, {BlackPawn, H3},
				{WhitePawn, A2}, {WhitePawn, B2}, {WhitePawn, C2}, {WhiteBishop, D2}, {WhiteBishop, E2}, {WhitePawn, F2}, {WhitePawn, G2}, {WhitePawn, H2},
				{WhiteRook, A1}, {WhiteKing, E1}, {WhiteRook, H1},
			},
			sideToMove: White, castling: AllCastling, epSq: NoSquare, halfmove: 0, ply: 0,
		},
		// TODO: Add more tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("unexpected error parsing valid FEN %q: %v", tt.fen, err)
			}
			want := wantPosition(tt.placements, tt.sideToMove, tt.castling, tt.epSq, tt.halfmove, tt.ply)
			assertPositionEqual(t, got, want, tt.fen)
		})
	}
}

func TestParseFEN_RoundTrip(t *testing.T) {
	tests := []string{
		startingFEN,
		"8/8/8/8/8/8/8/8 w - - 0 1",
		"rnbqkb1r/pppp1ppp/5n2/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
		"K6k/8/8/8/8/8/8/n6B w - - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
	}

	for _, fen := range tests {
		t.Run(fen, func(t *testing.T) {
			pos, err := ParseFEN(fen)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := pos.FEN()
			if got != fen {
				t.Fatalf("round-trip mismatch: started %q, got %q", fen, got)
			}
		})
	}
}

func TestParseFEN_FieldCount(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		wantErr bool
	}{
		{"empty string", "", true},
		{"only whitespace", "   ", true},
		{"missing all fields but placement", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR", true},
		{"missing halfmove and fullmove", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -", false},
		{"missing fullmove only", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0", false},
		{"too many fields", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 extra", true},
		{"extra trailing space", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 ", false},
		{"double space between fields", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w  KQkq - 0 1", false},
		{"kiwipete", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFEN(tt.fen)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error for %q but got none", tt.fen)
			}
		})
	}
}

func TestParseFEN_PiecePlacement(t *testing.T) {
	tests := []struct {
		name    string
		fen     string
		wantErr error
	}{
		{"too few ranks", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP w KQkq - 0 1", ErrInvalidRankCount},
		{"too many ranks", "rnbqkbnr/pppppppp/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidRankCount},
		{"rank too short", "rnbqkbn/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidRankLength},
		{"rank too long", "rnbqkbnrr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidRankLength},
		{"rank too long, digit overflow", "rnbqkbn9/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidRankDigit},
		{"rank sums over 8 with pieces and digit", "rnbqkbnr1/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidRankLength},
		{"digit zero forbidden", "rnbqkbnr/pppppppp/0/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidRankDigit},
		{"digit 9 invalid", "9/8/8/8/8/8/8/8 w - - 0 1", ErrInvalidRankDigit},
		{"invalid piece letter", "xnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidPieceChar},
		{"empty rank string", "//8/8/8/8/8/8 w - - 0 1", ErrInvalidRankLength},
		{"trailing slash", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR/ w KQkq - 0 1", ErrInvalidRankCount},
		{"leading slash", "/rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", ErrInvalidRankCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFEN(tt.fen)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v but got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPlyFullmovesRoundTrip(t *testing.T) {
	for ply := uint16(0); ply < 600; ply++ {
		fullmoves, side := PlyToFullmoves(ply)
		gotPly := FullmovesToPly(fullmoves, side)

		if gotPly != ply {
			t.Fatalf("ply=%d => fullmoves=%d side=%v => gotPly=%d",
				ply, fullmoves, side, gotPly)
		}
	}
}

func TestFullmovesPlyRoundTrip(t *testing.T) {
	t.Parallel()

	// fullmoves starts at 1
	for fm := uint16(1); fm < 300; fm++ {
		for _, side := range []Color{White, Black} {
			ply := FullmovesToPly(fm, side)
			gotFm, gotSide := PlyToFullmoves(ply)

			if gotFm != fm || gotSide != side {
				t.Fatalf("fullmoves round-trip failed: fullmoves=%d side=%v ply=%d gotFullmoves=%d gotSide=%v",
					fm, side, ply, gotFm, gotSide)
			}
		}
	}
}

func TestStartingPosition(t *testing.T) {
	want, err := ParseFEN(startingFEN)
	if err != nil {
		t.Fatalf("ParseFEN(startingFEN) returned error: %v", err)
	}

	got := StartingPosition()
	if got != want {
		t.Fatalf("StartingPosition mismatch.\nwant:\n%v\ngot:\n%v", want, got)
	}
}
