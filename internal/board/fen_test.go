package board

import "testing"

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
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
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
			sideToMove: White, castling: AllCastling, epSq: NoSquare, halfmove: 0, ply: 1,
		},
		{
			name:       "empty board",
			fen:        "8/8/8/8/8/8/8/8 w - - 0 1",
			placements: nil,
			sideToMove: White, castling: NoCastling, epSq: NoSquare, halfmove: 0, ply: 1,
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
			sideToMove: White, castling: AllCastling, epSq: E6, halfmove: 0, ply: 3,
		},

		{
			name:       "single piece each corner",
			fen:        "K6k/8/8/8/8/8/8/n6B w - - 0 1",
			placements: []placement{{WhiteKing, A8}, {BlackKing, H8}, {BlackKnight, A1}, {WhiteBishop, H1}},
			sideToMove: White, castling: NoCastling, epSq: NoSquare, halfmove: 0, ply: 1,
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
