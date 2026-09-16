package board

import (
	"errors"
	"testing"
)

func TestHas(t *testing.T) {
	tests := []struct {
		name  string
		c     Castling
		right CastlingRights
		want  bool
	}{
		{"has black kingside", Castling{NoSquare, H8, NoSquare, NoSquare}, BlackKingside, true},
		{"missing black queenside", Castling{NoSquare, H8, NoSquare, NoSquare}, BlackQueenside, false},
		{"has multiple rights", Castling{A8, H8, A1, H1}, WhiteQueenside, true},
		{"subset present", Castling{NoSquare, H8, NoSquare, H1}, BlackKingside, true},
		{"subset absent", Castling{NoSquare, H8, NoSquare, H1}, BlackQueenside, false},
		{"no castling has nothing", noCastling, BlackKingside, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Has(tt.right); got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestCastlingString(t *testing.T) {
	tests := []struct {
		name string
		c    Castling
		want string
	}{
		{"none", noCastling, "-"},
		{"black kingside", Castling{NoSquare, H8, NoSquare, NoSquare}, "k"},
		{"black queenside", Castling{A8, NoSquare, NoSquare, NoSquare}, "q"},
		{"white kingside", Castling{NoSquare, NoSquare, NoSquare, H1}, "K"},
		{"white queenside", Castling{NoSquare, NoSquare, A1, NoSquare}, "Q"},

		{"white kingside + queenside", Castling{NoSquare, NoSquare, A1, H1}, "KQ"},
		{"white kingside + black kingside", Castling{NoSquare, H8, NoSquare, H1}, "Kk"},
		{"white queenside + black kingside", Castling{NoSquare, H8, A1, NoSquare}, "Qk"},

		{"white all + black queenside", Castling{A8, NoSquare, A1, H1}, "KQq"},
		{"white queenside + black all", Castling{A8, H8, A1, NoSquare}, "Qkq"},

		{"all rights", Castling{A8, H8, A1, H1}, "KQkq"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(false); got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestParseCastlingRights(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Castling
		wantErr error
	}{
		{"none", "-", noCastling, nil},
		{"white kingside", "K", Castling{NoSquare, NoSquare, NoSquare, H1}, nil},
		{"white queenside", "Q", Castling{NoSquare, NoSquare, A1, NoSquare}, nil},
		{"black kingside", "k", Castling{NoSquare, H8, NoSquare, NoSquare}, nil},
		{"black queenside", "q", Castling{A8, NoSquare, NoSquare, NoSquare}, nil},

		{"white both", "KQ", Castling{NoSquare, NoSquare, A1, H1}, nil},
		{"black both", "kq", Castling{A8, H8, NoSquare, NoSquare}, nil},
		{"mixed all", "KQkq", Castling{A8, H8, A1, H1}, nil},
		{"mixed unordered", "qKkQ", Castling{A8, H8, A1, H1}, nil},

		{"empty string", "", noCastling, errInvalidCastlingLength},
		{"too long", "KQkq-", noCastling, errInvalidCastlingLength},
		{"invalid none and black kingside", "-k", noCastling, errInvalidCastlingChar},
		{"invalid black kingside and none", "k-", noCastling, errInvalidCastlingChar},

		{"invalid char letter", "X", noCastling, errInvalidCastlingChar},
		{"invalid char digit", "1", noCastling, errInvalidCastlingChar},
		{"invalid char symbol", "?", noCastling, errInvalidCastlingChar},

		{"duplicate white king", "KK", noCastling, errDuplicateCastlingChar},
		{"duplicate white queen", "QQ", noCastling, errDuplicateCastlingChar},
		{"duplicate black king", "kk", noCastling, errDuplicateCastlingChar},
		{"duplicate black queen", "qq", noCastling, errDuplicateCastlingChar},
		{"duplicate mixed", "KQK", noCastling, errDuplicateCastlingChar},
		{"duplicate across order", "qkq", noCastling, errDuplicateCastlingChar},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCastlingRights(tt.input, E1, E8)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v but got %v", tt.wantErr, err)
				}
				if got != noCastling {
					t.Errorf("expected zero Castling on error but got %v", got)
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
	withChess960(t)
	tests := []struct {
		name        string
		s           string
		whiteKingSq Square
		blackKingSq Square
		want        Castling
	}{
		{"standard start via file letters", "AHah", E1, E8, Castling{A8, H8, A1, H1}},
		{"white kingside only, rook on h", "H", E1, E8, Castling{NoSquare, NoSquare, NoSquare, H1}},
		{"white queenside only, rook on a", "A", E1, E8, Castling{NoSquare, NoSquare, A1, NoSquare}},
		{"black kingside only, rook on h", "h", E1, E8, Castling{NoSquare, H8, NoSquare, NoSquare}},
		{"black queenside only, rook on a", "a", E1, E8, Castling{A8, NoSquare, NoSquare, NoSquare}},
		{"960 king on b file, rooks on a and g", "GAga", B1, B8, Castling{A8, G8, A1, G1}},
		{"960 king on g file, rooks on h and d", "HDhd", G1, G8, Castling{D8, H8, D1, H1}},
		{"mixed case single rights", "Ab", E1, E8, Castling{B8, NoSquare, A1, NoSquare}},
		{"none", "-", E1, E8, noCastling},

		{"white both sides", "BD", C1, C8, Castling{NoSquare, NoSquare, B1, D1}},
		{"black both sides", "bd", C1, C8, Castling{B8, D8, NoSquare, NoSquare}},
		{"white kingside + black kingside", "Hh", E1, E8, Castling{NoSquare, H8, NoSquare, H1}},
		{"white queenside + black queenside", "Aa", E1, E8, Castling{A8, NoSquare, A1, NoSquare}},
		{"white kingside + black queenside", "Ha", E1, E8, Castling{A8, NoSquare, NoSquare, H1}},
		{"white queenside + black kingside", "Ah", E1, E8, Castling{NoSquare, H8, A1, NoSquare}},
		{"all four, mixed order", "haHA", E1, E8, Castling{A8, H8, A1, H1}},
		{"all four 960", "CHcf", G1, D8, Castling{C8, F8, C1, H1}},
		{"white both + black kingside", "AHh", E1, E8, Castling{NoSquare, H8, A1, H1}},
		{"white both + black queenside", "AHa", E1, E8, Castling{A8, NoSquare, A1, H1}},
		{"black both + white kingside", "ahH", E1, E8, Castling{A8, H8, NoSquare, H1}},
		{"black both + white queenside", "ahA", E1, E8, Castling{A8, H8, A1, NoSquare}},
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

func TestParseCastlingRights_UsesRealKingSquare(t *testing.T) {
	got, err := ParseCastlingRights("D", B1, B8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[WhiteKingside] != D1 {
		t.Errorf("king on B1: got %v, want WhiteKingside", got)
	}

	got, err = ParseCastlingRights("D", E1, E8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[WhiteQueenside] != D1 {
		t.Errorf("king on E1: got %v, want WhiteQueenside", got)
	}
}

func TestParseFEN_Chess960_NonEFileKing(t *testing.T) {
	withChess960(t)
	tests := []string{
		"1qrkrbbn/pppppppp/8/8/8/8/PPPPPPPP/1QRKRBBN w CEce - 0 1",
		"qnnbbrkr/pppppppp/8/8/8/8/PPPPPPPP/QNNBBRKR w FHfh - 0 1",
	}

	for _, fen := range tests {
		t.Run(fen, func(t *testing.T) {
			pos, err := ParseFEN(fen)
			if err != nil {
				t.Fatalf("ParseFEN failed: %v", err)
			}
			if got := pos.FEN(); got != fen {
				t.Errorf("got %s want %s", got, fen)
			}
		})
	}
}

func TestCastlingRoundTrip(t *testing.T) {
	tests := []string{
		"-",
		"K",
		"Q",
		"k",
		"q",
		"KQ",
		"kq",
		"Kk",
		"Kq",
		"Qk",
		"Qq",
		"KQk",
		"KQq",
		"Kkq",
		"Qkq",
		"KQkq",
	}

	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			got, err := ParseCastlingRights(s, E1, E8)
			if err != nil {
				t.Fatalf("unexpected error parsing %q: %v", s, err)
			}

			if got.String(false) != s {
				t.Fatalf("round-trip mismatch: started with %q, got %q", s, got.String(false))
			}
		})
	}
}

func TestCastlingRooks(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		want Castling
	}{
		{
			name: "standard starting position, all four rights",
			fen:  StartingFEN,
			want: Castling{
				WhiteKingside:  H1,
				WhiteQueenside: A1,
				BlackKingside:  H8,
				BlackQueenside: A8,
			},
		},
		{
			name: "no rights at all",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R w - - 0 1",
			want: Castling{
				WhiteKingside:  NoSquare,
				WhiteQueenside: NoSquare,
				BlackKingside:  NoSquare,
				BlackQueenside: NoSquare,
			},
		},
		{
			name: "only white kingside right, other rooks present but irrelevant",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R w K - 0 1",
			want: Castling{
				WhiteKingside:  H1,
				WhiteQueenside: NoSquare,
				BlackKingside:  NoSquare,
				BlackQueenside: NoSquare,
			},
		},
		{
			name: "only black queenside right",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R b q - 0 1",
			want: Castling{
				WhiteKingside:  NoSquare,
				WhiteQueenside: NoSquare,
				BlackKingside:  NoSquare,
				BlackQueenside: A8,
			},
		},
		{
			name: "960 white king on b file, white rooks on a and g, black king on d file, black rooks on c and f",
			fen:  "2rk1r2/8/8/8/8/8/8/RK4R1 w AGcf - 0 1",
			want: Castling{
				WhiteKingside:  G1,
				WhiteQueenside: A1,
				BlackKingside:  F8,
				BlackQueenside: C8,
			},
		},
		{
			name: "960 king on g-file, rooks on d and h",
			fen:  "3r2kr/8/8/8/8/8/8/3R2KR w HDhd - 0 1",
			want: Castling{
				WhiteKingside:  H1,
				WhiteQueenside: D1,
				BlackKingside:  H8,
				BlackQueenside: D8,
			},
		},
		{
			name: "right held but no matching rook on board returns NoSquare",
			fen:  "4k3/8/8/8/8/8/8/4K3 w KQkq - 0 1",
			want: Castling{
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

			got := FindCastlingRooks(&pos)
			if got != tt.want {
				t.Errorf("\n%v got %v want %v", pos.String(), got, tt.want)
			}
		})
	}
}
