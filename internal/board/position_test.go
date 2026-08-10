package board

import (
	"errors"
	"testing"
)

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
		{"has multiple rights", AllCastling, WhiteQueenside, true},
		{"subset present", BlackKingside | WhiteKingside, BlackKingside, true},
		{"subset absent", BlackKingside | WhiteKingside, BlackQueenside, false},
		{"no castling has nothing", NoCastling, BlackKingside, false},
		{"all castling has all", AllCastling, AllCastling, true},
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

		{"white kingside + queenside", WhiteKingside | WhiteQueenside, "KQ"},
		{"white kingside + black kingside", WhiteKingside | BlackKingside, "Kk"},
		{"white queenside + black kingside", WhiteQueenside | BlackKingside, "Qk"},

		{"white all + black queenside", WhiteKingside | WhiteQueenside | BlackQueenside, "KQq"},
		{"white queenside + black all", WhiteQueenside | BlackKingside | BlackQueenside, "Qkq"},

		{"all rights", AllCastling, "KQkq"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cr.String(); got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestParseCastlingRights(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    CastlingRights
		wantErr error
	}{
		{"none", "-", NoCastling, nil},
		{"white kingside", "K", WhiteKingside, nil},
		{"white queenside", "Q", WhiteQueenside, nil},
		{"black kingside", "k", BlackKingside, nil},
		{"black queenside", "q", BlackQueenside, nil},

		{"white both", "KQ", WhiteKingside | WhiteQueenside, nil},
		{"black both", "kq", BlackKingside | BlackQueenside, nil},
		{"mixed all", "KQkq", AllCastling, nil},
		{"mixed unordered", "qKkQ", AllCastling, nil},

		{"empty string", "", NoCastling, ErrInvalidCastlingLength},
		{"too long", "KQkq-", NoCastling, ErrInvalidCastlingLength},
		{"invalid none and black kingside", "-k", NoCastling, ErrInvalidCastlingChar},
		{"invalid black kingside and none", "k-", NoCastling, ErrInvalidCastlingChar},

		{"invalid char letter", "X", NoCastling, ErrInvalidCastlingChar},
		{"invalid char digit", "1", NoCastling, ErrInvalidCastlingChar},
		{"invalid char symbol", "?", NoCastling, ErrInvalidCastlingChar},

		{"duplicate white king", "KK", NoCastling, ErrDuplicateCastlingChar},
		{"duplicate white queen", "QQ", NoCastling, ErrDuplicateCastlingChar},
		{"duplicate black king", "kk", NoCastling, ErrDuplicateCastlingChar},
		{"duplicate black queen", "qq", NoCastling, ErrDuplicateCastlingChar},
		{"duplicate mixed", "KQK", NoCastling, ErrDuplicateCastlingChar},
		{"duplicate across order", "qkq", NoCastling, ErrDuplicateCastlingChar},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCastlingRights(tt.input, E1, E8)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v but got %v", tt.wantErr, err)
				}
				if got != NoCastling {
					t.Errorf("expected NoCastling on error but got %v", got)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestParseCastlingRights_Shredder(t *testing.T) {
	tests := []struct {
		name        string
		s           string
		whiteKingSq Square
		blackKingSq Square
		want        CastlingRights
	}{
		{"standard start via file letters", "AHah", E1, E8, AllCastling},
		{"white kingside only, rook on h", "H", E1, E8, WhiteKingside},
		{"white queenside only, rook on a", "A", E1, E8, WhiteQueenside},
		{"black kingside only, rook on h", "h", E1, E8, BlackKingside},
		{"black queenside only, rook on a", "a", E1, E8, BlackQueenside},
		{"960 king on b file, rooks on a and g", "GAga", B1, B8, AllCastling},
		{"960 king on g file, rooks on h and d", "HDhd", G1, G8, AllCastling},
		{"mixed case single rights", "Ab", E1, E8, WhiteQueenside | BlackQueenside},
		{"none", "-", E1, E8, NoCastling},

		{"white both sides", "BD", C1, C8, WhiteKingside | WhiteQueenside},
		{"black both sides", "bd", C1, C8, BlackKingside | BlackQueenside},
		{"white kingside + black kingside", "Hh", E1, E8, WhiteKingside | BlackKingside},
		{"white queenside + black queenside", "Aa", E1, E8, WhiteQueenside | BlackQueenside},
		{"white kingside + black queenside", "Ha", E1, E8, WhiteKingside | BlackQueenside},
		{"white queenside + black kingside", "Ah", E1, E8, WhiteQueenside | BlackKingside},
		{"all four, mixed order", "haHA", E1, E8, AllCastling},
		{"all four 960", "CHcf", G1, D8, AllCastling},
		{"white both + black kingside", "AHh", E1, E8, WhiteKingside | WhiteQueenside | BlackKingside},
		{"white both + black queenside", "AHa", E1, E8, WhiteKingside | WhiteQueenside | BlackQueenside},
		{"black both + white kingside", "ahH", E1, E8, BlackKingside | BlackQueenside | WhiteKingside},
		{"black both + white queenside", "ahA", E1, E8, BlackKingside | BlackQueenside | WhiteQueenside},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCastlingRights(tt.s, tt.whiteKingSq, tt.blackKingSq)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestParseCastlingRights_ShredderErrors(t *testing.T) {
	tests := []struct {
		name     string
		from, to Square
		s        string
	}{
		{"empty string", E1, E8, ""},
		{"too long", D1, E8, "AHahb"},
		{"duplicate file letter", C1, H8, "AA"},
		{"file letter equal to white king's file", E1, A8, "E"},
		{"file letter equal to black king's file", H1, F8, "Dbf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCastlingRights(tt.s, tt.from, tt.to)
			if err == nil {
				t.Error("expected an error, got none")
			}
		})
	}
}

func TestCastlingRightsRoundTrip(t *testing.T) {
	tests := []CastlingRights{
		NoCastling,
		BlackKingside,
		BlackQueenside,
		WhiteKingside,
		WhiteQueenside,
		BlackKingside | BlackQueenside,
		WhiteKingside | WhiteQueenside,
		WhiteKingside | BlackKingside,
		WhiteKingside | BlackQueenside,
		WhiteQueenside | BlackKingside,
		WhiteQueenside | BlackQueenside,
		WhiteKingside | WhiteQueenside | BlackKingside,
		WhiteKingside | WhiteQueenside | BlackQueenside,
		WhiteKingside | BlackKingside | BlackQueenside,
		WhiteQueenside | BlackKingside | BlackQueenside,
		AllCastling,
	}

	for _, cr := range tests {
		t.Run(cr.String(), func(t *testing.T) {
			s := cr.String()

			got, err := ParseCastlingRights(s, E1, E8)
			if err != nil {
				t.Fatalf("unexpected error parsing %q: %v", s, err)
			}
			if got != cr {
				t.Fatalf("round-trip mismatch: started with %v, string %q, parsed %v", cr, s, got)
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
		{"remove black kingside", AllCastling, BlackKingside, BlackQueenside | WhiteKingside | WhiteQueenside},
		{"remove black queenside", AllCastling, BlackQueenside, BlackKingside | WhiteKingside | WhiteQueenside},
		{"remove white kingside", AllCastling, WhiteKingside, BlackKingside | BlackQueenside | WhiteQueenside},
		{"remove white queenside", AllCastling, WhiteQueenside, BlackKingside | BlackQueenside | WhiteKingside},
		{"remove absent right", BlackKingside, BlackQueenside, BlackKingside},
		{"remove only right", WhiteQueenside, WhiteQueenside, NoCastling},
		{"remove from none", NoCastling, BlackKingside, NoCastling},
		{"remove all rights", AllCastling, AllCastling, NoCastling},
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
	pos := StartingPosition()

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {

		for sq := range H8 {
			pos.RemovePiece(sq)
		}
	}
}

func TestPieceOn(t *testing.T) {
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
	pos := StartingPosition()

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

func TestCastlingRooks(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want CastlingRookSquares
	}{
		{
			name: "standard starting position, all four rights",
			fen:  startingFEN,
			want: CastlingRookSquares{
				WhiteKingside:  H1,
				WhiteQueenside: A1,
				BlackKingside:  H8,
				BlackQueenside: A8,
			},
		},
		{
			name: "no rights at all",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R w - - 0 1",
			want: CastlingRookSquares{
				WhiteKingside:  NoSquare,
				WhiteQueenside: NoSquare,
				BlackKingside:  NoSquare,
				BlackQueenside: NoSquare,
			},
		},
		{
			name: "only white kingside right, other rooks present but irrelevant",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R w K - 0 1",
			want: CastlingRookSquares{
				WhiteKingside:  H1,
				WhiteQueenside: NoSquare,
				BlackKingside:  NoSquare,
				BlackQueenside: NoSquare,
			},
		},
		{
			name: "only black queenside right",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R b q - 0 1",
			want: CastlingRookSquares{
				WhiteKingside:  NoSquare,
				WhiteQueenside: NoSquare,
				BlackKingside:  NoSquare,
				BlackQueenside: A8,
			},
		},
		{
			name: "960 white king on b file, white rooks on a and g, black king on d file, black rooks on c and f",
			fen:  "2rk1r2/8/8/8/8/8/8/RK4R1 w AGcf - 0 1",
			want: CastlingRookSquares{
				WhiteKingside:  G1,
				WhiteQueenside: A1,
				BlackKingside:  F8,
				BlackQueenside: C8,
			},
		},
		{
			name: "960 king on g-file, rooks on d and h",
			fen:  "3r2kr/8/8/8/8/8/8/3R2KR w HDhd - 0 1",
			want: CastlingRookSquares{
				WhiteKingside:  H1,
				WhiteQueenside: D1,
				BlackKingside:  H8,
				BlackQueenside: D8,
			},
		},
		{
			name: "right held but no matching rook on board (malformed/edge state) returns NoSquare",
			fen:  "4k3/8/8/8/8/8/8/4K3 w KQkq - 0 1",
			want: CastlingRookSquares{
				WhiteKingside:  NoSquare,
				WhiteQueenside: NoSquare,
				BlackKingside:  NoSquare,
				BlackQueenside: NoSquare,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("bad test FEN: %v", err)
			}

			got := pos.CastlingRooks()
			if got != tt.want {
				t.Errorf("\n%v got %v want %v", pos, got, tt.want)
			}
		})
	}
}
