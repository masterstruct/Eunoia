package board

import "testing"

type placement struct {
	p  Piece
	sq Square
}

func TestHas(t *testing.T) {
	tests := []struct {
		name  string
		cr    CastlingRights
		right CastlingRights
		want  bool
	}{
		{"has black kingside", BlackKingside, BlackKingside, true},
		{"missing black queenside", BlackKingside, BlackQueenside, false},
		{"has multiple rights", AnyCastling, WhiteQueenside, true},
		{"subset present", BlackKingside | WhiteKingside, BlackKingside, true},
		{"subset absent", BlackKingside | WhiteKingside, BlackQueenside, false},
		{"no castling has nothing", NoCastling, BlackKingside, false},
		{"any castling has any", AnyCastling, AnyCastling, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cr.Has(tt.right); got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestCastlingRightsString(t *testing.T) {
	tests := []struct {
		name string
		cr   CastlingRights
		want string
	}{
		{"none", NoCastling, "-"},
		{"black kingside", BlackKingside, "k"},
		{"black queenside", BlackQueenside, "q"},
		{"white kingside", WhiteKingside, "K"},
		{"white queenside", WhiteQueenside, "Q"},

		{"black kingside + queenside", BlackKingside | BlackQueenside, "kq"},
		{"white kingside + queenside", WhiteKingside | WhiteQueenside, "KQ"},
		{"white kingside + black kingside", WhiteKingside | BlackKingside, "Kk"},
		{"white kingside + black queenside", WhiteKingside | BlackQueenside, "Kq"},
		{"white queenside + black kingside", WhiteQueenside | BlackKingside, "Qk"},
		{"white queenside + black queenside", WhiteQueenside | BlackQueenside, "Qq"},

		{"white all + black kingside", WhiteKingside | WhiteQueenside | BlackKingside, "KQk"},
		{"white all + black queenside", WhiteKingside | WhiteQueenside | BlackQueenside, "KQq"},
		{"white kingside + black all", WhiteKingside | BlackKingside | BlackQueenside, "Kkq"},
		{"white queenside + black all", WhiteQueenside | BlackKingside | BlackQueenside, "Qkq"},

		{"all rights", AnyCastling, "KQkq"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cr.String(); got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name  string
		cr    CastlingRights
		right CastlingRights
		want  CastlingRights
	}{
		{"remove black kingside", AnyCastling, BlackKingside, BlackQueenside | WhiteKingside | WhiteQueenside},
		{"remove black queenside", AnyCastling, BlackQueenside, BlackKingside | WhiteKingside | WhiteQueenside},
		{"remove white kingside", AnyCastling, WhiteKingside, BlackKingside | BlackQueenside | WhiteQueenside},
		{"remove white queenside", AnyCastling, WhiteQueenside, BlackKingside | BlackQueenside | WhiteKingside},
		{"remove absent right", BlackKingside, BlackQueenside, BlackKingside},
		{"remove only right", WhiteQueenside, WhiteQueenside, NoCastling},
		{"remove from none", NoCastling, BlackKingside, NoCastling},
		{"remove all rights", AnyCastling, AnyCastling, NoCastling},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := tt.cr
			cr.Remove(tt.right)

			if cr != tt.want {
				t.Errorf("removed %v but got %v, want %v", tt.right, cr, tt.want)
			}
		})
	}
}

func TestPieceBB(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name       string
		placements []placement
		piece      Piece
		want       Bitboard
	}{
		{
			name:       "white queen present",
			placements: []placement{{WhiteQueen, D1}},
			piece:      WhiteQueen,
			want:       setBits([]Square{D1}),
		},
		{
			name:       "wrong color",
			placements: []placement{{BlackQueen, D1}},
			piece:      WhiteQueen,
			want:       EmptyBB,
		},
		{
			name:       "wrong piece type",
			placements: []placement{{WhiteRook, D1}},
			piece:      WhiteQueen,
			want:       EmptyBB,
		},
		{
			name:       "mixed board",
			placements: []placement{{WhiteQueen, D1}, {BlackQueen, D8}},
			piece:      WhiteQueen,
			want:       setBits([]Square{D1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition()
			for _, pl := range tt.placements {
				pos.PlacePiece(pl.p, pl.sq)
			}

			got := pos.PieceBB(tt.piece)
			if got != tt.want {
				t.Fatalf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestOccupied(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name  string
		place []placement
		want  Bitboard
	}{
		{
			name: "empty",
			want: EmptyBB,
		},
		{
			name: "single piece",
			place: []placement{
				{WhiteQueen, D1},
			},
			want: setBits([]Square{D1}),
		},
		{
			name: "multiple pieces different colors",
			place: []placement{
				{WhiteQueen, D1},
				{BlackKing, E8},
				{WhiteRook, A1},
			},
			want: setBits([]Square{D1, E8, A1}),
		},
		{
			name: "duplicate same square",
			place: []placement{
				{WhiteQueen, D1},
				{WhiteQueen, D1},
			},
			want: setBits([]Square{D1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition()
			for _, pl := range tt.place {
				pos.PlacePiece(pl.p, pl.sq)
			}
			got := pos.Occupied()
			if got != tt.want {
				t.Fatalf("expected %v but got %v", tt.want, got)
			}
			if pos.Occupied() != allPieceBB(pos) {
				t.Fatalf("occupied and piece bitboards out of sync: %v vs %v", pos.Occupied(), allPieceBB(pos))
			}
		})
	}
}

func TestPlacePiece(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name       string
		placements []placement
		pieces     map[Piece]Bitboard
		want       Bitboard
	}{
		{
			name:       "single white queen",
			placements: []placement{{WhiteQueen, D1}},
			pieces: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
			},
			want: setBits([]Square{D1}),
		},
		{
			name:       "single black king",
			placements: []placement{{BlackKing, E8}},
			pieces: map[Piece]Bitboard{
				BlackKing: setBits([]Square{E8}),
			},
			want: setBits([]Square{E8}),
		},
		{
			name: "multiple pieces different colors",
			placements: []placement{
				{WhiteQueen, D1},
				{BlackKing, E8},
				{WhiteRook, A1},
			},
			pieces: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
				BlackKing:  setBits([]Square{E8}),
				WhiteRook:  setBits([]Square{A1}),
			},
			want: setBits([]Square{D1, E8, A1}),
		},
		{
			name: "same square different pieces",
			placements: []placement{
				{WhiteQueen, D1},
				{WhiteRook, D1},
			},
			pieces: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
				WhiteRook:  setBits([]Square{D1}),
			},
			want: setBits([]Square{D1}),
		},
		{
			name: "duplicate same piece same square",
			placements: []placement{
				{WhiteQueen, D1},
				{WhiteQueen, D1},
			},
			pieces: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
			},
			want: setBits([]Square{D1}),
		},
		{
			name: "multiple same color pieces",
			placements: []placement{
				{WhiteKing, E1},
				{WhiteQueen, D1},
				{WhiteRook, A1},
			},
			pieces: map[Piece]Bitboard{
				WhiteKing:  setBits([]Square{E1}),
				WhiteQueen: setBits([]Square{D1}),
				WhiteRook:  setBits([]Square{A1}),
			},
			want: setBits([]Square{E1, D1, A1}),
		},
		{
			name: "board with all piece types",
			placements: []placement{
				{WhitePawn, A2},
				{WhiteKnight, B1},
				{WhiteBishop, C1},
				{WhiteRook, D1},
				{WhiteQueen, E1},
				{WhiteKing, F1},
				{BlackPawn, A7},
				{BlackKnight, B8},
				{BlackBishop, C8},
				{BlackRook, D8},
				{BlackQueen, E8},
				{BlackKing, F8},
			},
			pieces: map[Piece]Bitboard{
				WhitePawn:   setBits([]Square{A2}),
				WhiteKnight: setBits([]Square{B1}),
				WhiteBishop: setBits([]Square{C1}),
				WhiteRook:   setBits([]Square{D1}),
				WhiteQueen:  setBits([]Square{E1}),
				WhiteKing:   setBits([]Square{F1}),
				BlackPawn:   setBits([]Square{A7}),
				BlackKnight: setBits([]Square{B8}),
				BlackBishop: setBits([]Square{C8}),
				BlackRook:   setBits([]Square{D8}),
				BlackQueen:  setBits([]Square{E8}),
				BlackKing:   setBits([]Square{F8}),
			},
			want: setBits([]Square{A2, B1, C1, D1, E1, F1, A7, B8, C8, D8, E8, F8}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition()

			for _, pl := range tt.placements {
				pos.PlacePiece(pl.p, pl.sq)

				if pos.Board[pl.sq] != pl.p {
					t.Fatalf("placed piece %v on square %v but board did not update", pl.p, pl.sq)
				}
			}

			if got := pos.Occupied(); got != tt.want {
				t.Fatalf("occupied: expected %v but got %v", tt.want, got)
			}

			if got := allPieceBB(pos); got != tt.want {
				t.Fatalf("pieces: expected %v but got %v", tt.want, got)
			}

			for piece, want := range tt.pieces {
				if got := pos.PieceBB(piece); got != want {
					t.Fatalf("expected %v but got %v", want, got)
				}
			}
		})
	}
}

func TestRemovePiece(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name   string
		place  []placement
		remove Square
		want   Bitboard
	}{
		{
			name: "empty",
			want: EmptyBB,
		},
		{
			name: "remove single piece",
			place: []placement{
				{WhiteQueen, D1},
			},
			remove: D1,
			want:   EmptyBB,
		},
		{
			name: "remove one of multiple pieces",
			place: []placement{
				{WhiteQueen, D1},
				{BlackKnight, E8},
				{WhiteRook, A1},
			},
			remove: E8,
			want:   setBits([]Square{D1, A1}),
		},
		{
			name: "remove nonexistent piece leaves board unchanged",
			place: []placement{
				{WhiteQueen, D1},
			},
			remove: E1,
			want:   setBits([]Square{D1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition()
			for _, pl := range tt.place {
				pos.PlacePiece(pl.p, pl.sq)
			}

			pos.RemovePiece(tt.remove)

			if got := pos.Occupied(); got != tt.want {
				t.Fatalf("expected %v but got %v", tt.want, got)
			}

			if got := allPieceBB(pos); got != tt.want {
				t.Fatalf("expected %v but got %v", tt.want, got)
			}

			for _, pl := range tt.place {
				if tt.remove == pl.sq && pos.Board[pl.sq] != NoPiece {
					t.Fatalf("removed piece %v on square %v but board did not update", pl.p, pl.sq)
				}
			}
		})
	}
}

func BenchmarkRemovePiece(b *testing.B) {
	InitBitboards()

	pos := NewPosition()

	// TODO: replace with FEN loading
	setup := []placement{
		{BlackRook, A8}, {BlackBishop, B8}, {BlackKnight, C8}, {BlackQueen, D8},
		{BlackKing, E8}, {BlackBishop, F8}, {BlackKnight, G8}, {BlackRook, H8},
		{BlackPawn, A7}, {BlackPawn, B7}, {BlackPawn, C7}, {BlackPawn, D7},
		{BlackPawn, E7}, {BlackPawn, F7}, {BlackPawn, G7}, {BlackPawn, H7},

		{WhitePawn, A2}, {WhitePawn, B2}, {WhitePawn, C2}, {WhitePawn, D2},
		{WhitePawn, E2}, {WhitePawn, F2}, {WhitePawn, G2}, {WhitePawn, H2},
		{WhiteRook, A1}, {WhiteBishop, B1}, {WhiteKnight, C1}, {WhiteQueen, D1},
		{WhiteKing, E1}, {WhiteBishop, F1}, {WhiteKnight, G1}, {WhiteRook, H1},
	}

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		for _, x := range setup {
			pos.PlacePiece(x.p, x.sq)
		}

		for sq := range H8 {
			pos.RemovePiece(sq)
		}
	}
}

func TestPieceOn(t *testing.T) {
	InitBitboards()

	pos := NewPosition()
	pos.PlacePiece(WhiteKing, E1)
	pos.PlacePiece(BlackPawn, E7)

	tests := []struct {
		name   string
		sq     Square
		want   Piece
		wantOk bool
	}{
		{"white king on e1", E1, WhiteKing, true},
		{"black pawn on e7", E7, BlackPawn, true},
		{"empty square", E4, NoPiece, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := pos.PieceOn(tt.sq)
			if ok != tt.wantOk {
				t.Errorf("expected ok=%v but got %v", tt.wantOk, ok)
			}
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func BenchmarkPieceOn(b *testing.B) {
	InitBitboards()

	pos := NewPosition()

	// TODO: replace with FEN loading
	setup := []placement{
		{BlackRook, A8}, {BlackBishop, B8}, {BlackKnight, C8}, {BlackQueen, D8},
		{BlackKing, E8}, {BlackBishop, F8}, {BlackKnight, G8}, {BlackRook, H8},
		{BlackPawn, A7}, {BlackPawn, B7}, {BlackPawn, C7}, {BlackPawn, D7},
		{BlackPawn, E7}, {BlackPawn, F7}, {BlackPawn, G7}, {BlackPawn, H7},

		{WhitePawn, A2}, {WhitePawn, B2}, {WhitePawn, C2}, {WhitePawn, D2},
		{WhitePawn, E2}, {WhitePawn, F2}, {WhitePawn, G2}, {WhitePawn, H2},
		{WhiteRook, A1}, {WhiteBishop, B1}, {WhiteKnight, C1}, {WhiteQueen, D1},
		{WhiteKing, E1}, {WhiteBishop, F1}, {WhiteKnight, G1}, {WhiteRook, H1},
	}

	for _, x := range setup {
		pos.PlacePiece(x.p, x.sq)
	}

	var piece Piece
	var ok bool
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		for sq := range A8 {
			piece, ok = pos.PieceOn(sq)
		}
	}
	_, _ = piece, ok
}

func setBits(sqs []Square) Bitboard {
	var bb Bitboard
	for _, sq := range sqs {
		bb.SetBit(sq)
	}
	return bb
}

func allPieceBB(pos Position) Bitboard {
	b := EmptyBB
	for _, pt := range PieceTypes() {
		b |= pos.pieces[pt]
	}
	return b
}

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	for sq, piece := range board {
		if piece != NoPiece {
			t.Errorf("square %d: expected %v but got %v", sq, NoPiece, piece)
		}
	}
}
