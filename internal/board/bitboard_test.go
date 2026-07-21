package board

import "testing"

func TestSetBit(t *testing.T) {
	InitBitboards()

	tests := []struct {
		name string
		bb   Bitboard
		sq   uint8
		want Bitboard
	}{
		{"set bit 0", EmptyBB, 0, 1 << 0},
		{"set bit 1", EmptyBB, 1, 1 << 1},
		{"set bit 3", EmptyBB, 3, 1 << 3},
		{"set bit 8", EmptyBB, 8, 1 << 8},
		{"set bit 9", EmptyBB, 9, 1 << 9},
		{"set bit 62", EmptyBB, 62, 1 << 62},
		{"set bit 63", EmptyBB, 63, 1 << 63},
		{"set on existing bitboard", 1, 1, 3},
		{"set bit 8 on existing", 1, 8, 257},
		{"set bit 63 on existing", 256, 63, (1 << 63) + 256},
		{"set already set bit", 256, 8, 256},
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
		sq   uint8
		want Bitboard
	}{
		{"clear bit 0", 1 << 0, 0, EmptyBB},
		{"clear bit 1", 1 << 1, 1, EmptyBB},
		{"clear bit 3", 1 << 3, 3, EmptyBB},
		{"clear bit 8", 1 << 8, 8, EmptyBB},
		{"clear bit 9", 1 << 9, 9, EmptyBB},
		{"clear bit 62", 1 << 62, 62, EmptyBB},
		{"clear bit 63", 1 << 63, 63, EmptyBB},
		{"clear one of many bits", (1 << 0) | (1 << 1), 1, 1 << 0},
		{"clear bit 8 from existing", (1 << 0) | (1 << 8), 8, 1 << 0},
		{"clear bit 63 from existing", (1 << 8) | (1 << 63), 63, 1 << 8},
		{"clear already clear bit", 256, 0, 256},
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
		sq   uint8
		want bool
	}{
		{"empty bitboard", EmptyBB, 0, false},
		{"bit 0 set", 1 << 0, 0, true},
		{"bit 1 set", 1 << 1, 1, true},
		{"bit 3 set", 1 << 3, 3, true},
		{"bit 8 set", 1 << 8, 8, true},
		{"bit 9 set", 1 << 9, 9, true},
		{"bit 62 set", 1 << 62, 62, true},
		{"bit 63 set", 1 << 63, 63, true},
		{"wrong bit not set", 1 << 3, 8, false},
		{"one of many bits set", (1 << 0) | (1 << 8), 8, true},
		{"bit not present in many", (1 << 0) | (1 << 8), 63, false},
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
		want int
	}{
		{"empty", EmptyBB, -1},
		{"bit 0", 1 << 0, 0},
		{"bit 1", 1 << 1, 1},
		{"bit 3", 1 << 3, 3},
		{"bit 8", 1 << 8, 8},
		{"bit 63", 1 << 63, 63},
		{"multiple bits low first", (1 << 3) | (1 << 8), 3},
		{"multiple bits high first", (1 << 0) | (1 << 63), 0},
		{"full board", FullBB, 0},
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
