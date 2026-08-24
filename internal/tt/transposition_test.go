package tt

import (
	"testing"
	"unsafe"
)

func TestNextPow2(t *testing.T) {
	tests := []struct {
		in   uint
		want uint
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{17, 32},
		{1000, 1024},
		{3498, 4096},
	}

	for _, tt := range tests {
		if got := nextPow2(tt.in); got != tt.want {
			t.Fatalf("%d: got %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestTTSizeFromMB(t *testing.T) {
	const entrySize = uint(unsafe.Sizeof(Entry{}))

	tests := []struct {
		name string
		mb   uint
		want uint
	}{
		{
			name: "0mb",
			mb:   0,
			want: 1,
		},
		{
			name: "1mb",
			mb:   1,
			want: nextPow2((1 * 1024 * 1024) / entrySize),
		},
		{
			name: "8mb",
			mb:   8,
			want: nextPow2((8 * 1024 * 1024) / entrySize),
		},
		{
			name: "16mb",
			mb:   16,
			want: nextPow2((16 * 1024 * 1024) / entrySize),
		},
		{
			name: "2gb",
			mb:   2 * 1024,
			want: nextPow2((2 * 1024 * 1024 * 1024) / entrySize),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SizeFromMB(tt.mb)
			if got != tt.want {
				t.Fatalf("%d: got %d, want %d", tt.mb, got, tt.want)
			}
		})
	}
}
