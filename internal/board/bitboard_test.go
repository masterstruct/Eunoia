package board

import "testing"

func TestSetBit(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name string
		bb   Bitboard
		sq   Square
		want Bitboard
	}{
		{"set bit 0", EmptyBB, A1, 1 << 0},
		{"set bit 1", EmptyBB, B1, 1 << 1},
		{"set bit 3", EmptyBB, D1, 1 << 3},
		{"set bit 8", EmptyBB, A2, 1 << 8},
		{"set bit 9", EmptyBB, B2, 1 << 9},
		{"set bit 62", EmptyBB, G8, 1 << 62},
		{"set bit 63", EmptyBB, H8, 1 << 63},
		{"set on existing bitboard", 1, B1, 3},
		{"set bit 8 on existing", 1, A2, 257},
		{"set bit 63 on existing", 256, H8, (1 << 63) + 256},
		{"set already set bit", 256, A2, 256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bb.SetBit(tt.sq)
			if tt.bb != tt.want {
				t.Errorf("expected %d but got %d", tt.want, tt.bb)
			}
		})
	}
}

func TestClearBit(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name string
		bb   Bitboard
		sq   Square
		want Bitboard
	}{
		{"clear bit 0", 1 << 0, A1, EmptyBB},
		{"clear bit 1", 1 << 1, B1, EmptyBB},
		{"clear bit 3", 1 << 3, D1, EmptyBB},
		{"clear bit 8", 1 << 8, A2, EmptyBB},
		{"clear bit 9", 1 << 9, B2, EmptyBB},
		{"clear bit 62", 1 << 62, G8, EmptyBB},
		{"clear bit 63", 1 << 63, H8, EmptyBB},
		{"clear one of many bits", (1 << 0) | (1 << 1), B1, 1 << 0},
		{"clear bit 8 from existing", (1 << 0) | (1 << 8), A2, 1 << 0},
		{"clear bit 63 from existing", (1 << 8) | (1 << 63), H8, 1 << 8},
		{"clear already clear bit", 256, A1, 256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bb.ClearBit(tt.sq)
			if tt.bb != tt.want {
				t.Errorf("expected %d but got %d", tt.want, tt.bb)
			}
		})
	}
}

func TestIsBitSet(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name string
		bb   Bitboard
		sq   Square
		want bool
	}{
		{"empty bitboard", EmptyBB, A1, false},
		{"bit 0 set", 1 << 0, A1, true},
		{"bit 1 set", 1 << 1, B1, true},
		{"bit 3 set", 1 << 3, D1, true},
		{"bit 8 set", 1 << 8, A2, true},
		{"bit 9 set", 1 << 9, B2, true},
		{"bit 62 set", 1 << 62, G8, true},
		{"bit 63 set", 1 << 63, H8, true},
		{"wrong bit not set", 1 << 3, A2, false},
		{"one of many bits set", (1 << 0) | (1 << 8), A2, true},
		{"bit not present in many", (1 << 0) | (1 << 8), H8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bb.IsBitSet(tt.sq)
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestCountBits(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name string
		bb   Bitboard
		want int
	}{
		{"empty bitboard", EmptyBB, 0},
		{"one bit set", 1 << 0, 1},
		{"bit 3 set", 1 << 3, 1},
		{"two bits set", (1 << 0) | (1 << 3), 2},
		{"three bits set", (1 << 0) | (1 << 3) | (1 << 8), 3},
		{"low bits set", 0b1111, 4},
		{"non-adjacent bits", (1 << 1) | (1 << 9) | (1 << 63), 3},
		{"half board", 0x00000000FFFFFFFF, 32},
		{"full board", FullBB, 64},
		{"single high bit", 1 << 63, 1},
		{"alternating bits", 0xAAAAAAAAAAAAAAAA, 32},
		{"other alternating bits", 0x5555555555555555, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bb.CountBits()
			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}

func TestLSB(t *testing.T) {
	tests := []struct {
		name string
		bb   Bitboard
		want Square
	}{
		{"empty", EmptyBB, NoSquare},
		{"bit 0", 1 << 0, A1},
		{"bit 1", 1 << 1, B1},
		{"bit 3", 1 << 3, D1},
		{"bit 8", 1 << 8, A2},
		{"bit 63", 1 << 63, H8},
		{"multiple bits low first", (1 << 3) | (1 << 8), D1},
		{"multiple bits high first", (1 << 0) | (1 << 63), A1},
		{"full board", FullBB, A1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bb.LSB()
			if got != tt.want {
				t.Errorf("expected %d but got %d", tt.want, got)
			}
		})
	}
}

func TestPopLSB(t *testing.T) {
	tests := []struct {
		name   string
		bb     Bitboard
		wantBB Bitboard
		wantSq Square
	}{
		{"empty", EmptyBB, EmptyBB, NoSquare},
		{"single bit 0", 1 << 0, EmptyBB, A1},
		{"single bit 3", 1 << 3, EmptyBB, D1},
		{"single high bit", 1 << 63, EmptyBB, H8},
		{"two bits low first", (1 << 3) | (1 << 8), 1 << 8, D1},
		{"two bits high first", (1 << 0) | (1 << 63), 1 << 63, A1},
		{"many bits", (1 << 0) | (1 << 3) | (1 << 8), (1 << 3) | (1 << 8), A1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := tt.bb
			gotSq := bb.PopLSB()

			if gotSq != tt.wantSq {
				t.Errorf("expected square %d but got %d", tt.wantSq, gotSq)
			}
			if bb != tt.wantBB {
				t.Errorf("expected bitboard %d but got %d", tt.wantBB, bb)
			}
		})
	}
}
