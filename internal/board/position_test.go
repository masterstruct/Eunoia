package board

import "testing"

type placement struct {
	p  Piece
	sq Square
}

func TestPieceBB(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name string
		pos  Position
		p    Piece
		want Bitboard
	}{
		{
			name: "white queen present",
			pos: func() Position {
				pos := NewPosition()
				pos.PlacePiece(WhiteQueen, D1)
				return pos
			}(),
			p:    WhiteQueen,
			want: setBits([]Square{D1}),
		},
		{
			name: "wrong color",
			pos: func() Position {
				pos := NewPosition()
				pos.PlacePiece(BlackQueen, D1)
				return pos
			}(),
			p:    WhiteQueen,
			want: EmptyBB,
		},
		{
			name: "wrong piece type",
			pos: func() Position {
				pos := NewPosition()
				pos.PlacePiece(WhiteRook, D1)
				return pos
			}(),
			p:    WhiteQueen,
			want: EmptyBB,
		},
		{
			name: "mixed board",
			pos: func() Position {
				pos := NewPosition()
				pos.PlacePiece(WhiteQueen, D1)
				pos.PlacePiece(BlackQueen, D8)
				return pos
			}(),
			p:    WhiteQueen,
			want: setBits([]Square{D1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pos.PieceBB(tt.p)
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
		name         string
		place        []placement
		wantOccupied Bitboard
		wantPiece    map[Piece]Bitboard
	}{
		{
			name:         "single white queen",
			place:        []placement{{WhiteQueen, D1}},
			wantOccupied: setBits([]Square{D1}),
			wantPiece: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
			},
		},
		{
			name:         "single black king",
			place:        []placement{{BlackKing, E8}},
			wantOccupied: setBits([]Square{E8}),
			wantPiece: map[Piece]Bitboard{
				BlackKing: setBits([]Square{E8}),
			},
		},
		{
			name: "multiple pieces different colors",
			place: []placement{
				{WhiteQueen, D1},
				{BlackKing, E8},
				{WhiteRook, A1},
			},
			wantOccupied: setBits([]Square{D1, E8, A1}),
			wantPiece: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
				BlackKing:  setBits([]Square{E8}),
				WhiteRook:  setBits([]Square{A1}),
			},
		},
		{
			name: "same square different pieces",
			place: []placement{
				{WhiteQueen, D1},
				{WhiteRook, D1},
			},
			wantOccupied: setBits([]Square{D1}),
			wantPiece: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
				WhiteRook:  setBits([]Square{D1}),
			},
		},
		{
			name: "duplicate same piece same square",
			place: []placement{
				{WhiteQueen, D1},
				{WhiteQueen, D1},
			},
			wantOccupied: setBits([]Square{D1}),
			wantPiece: map[Piece]Bitboard{
				WhiteQueen: setBits([]Square{D1}),
			},
		},
		{
			name: "multiple same color pieces",
			place: []placement{
				{WhiteKing, E1},
				{WhiteQueen, D1},
				{WhiteRook, A1},
			},
			wantOccupied: setBits([]Square{E1, D1, A1}),
			wantPiece: map[Piece]Bitboard{
				WhiteKing:  setBits([]Square{E1}),
				WhiteQueen: setBits([]Square{D1}),
				WhiteRook:  setBits([]Square{A1}),
			},
		},
		{
			name: "board with all piece types",
			place: []placement{
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
			wantOccupied: setBits([]Square{A2, B1, C1, D1, E1, F1, A7, B8, C8, D8, E8, F8}),
			wantPiece: map[Piece]Bitboard{
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := NewPosition()

			for _, pl := range tt.place {
				pos.PlacePiece(pl.p, pl.sq)
			}

			if got := pos.Occupied(); got != tt.wantOccupied {
				t.Fatalf("occupied: expected %v but got %v", tt.wantOccupied, got)
			}

			if got := allPieceBB(pos); got != tt.wantOccupied {
				t.Fatalf("pieces: expected %v but got %v", tt.wantOccupied, got)
			}

			for p, want := range tt.wantPiece {
				if got := pos.PieceBB(p); got != want {
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
		})
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

	for sq, p := range board {
		if p != NoPiece {
			t.Errorf("square %d: expected %v but got %v", sq, NoPiece, p)
		}
	}
}
