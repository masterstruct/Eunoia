package tt

import (
	"testing"
	"unsafe"

	"github.com/masterstruct/Eunoia/internal/board"
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

func TestStoreAndProbe(t *testing.T) {
	var tt Table
	var move board.Move
	key := uint64(42)

	tt.Store(key, move, 100, 5, Exact)

	entry, hit := tt.Probe(key)
	if !hit {
		t.Fatalf("expected hit for stored key")
	}
	if entry.Key != key || entry.Score != 100 || entry.Depth != 5 || entry.Flag != Exact {
		t.Errorf("got key=%d score=%d depth=%d flag=%v", entry.Key, entry.Score, entry.Depth, entry.Flag)
	}
}

func TestProbeMiss(t *testing.T) {
	var tt Table
	_, hit := tt.Probe(999999)
	if hit {
		t.Errorf("expected miss on empty table")
	}
}

func TestStoreOverwrites(t *testing.T) {
	var tt Table
	var move board.Move
	key := uint64(7)

	tt.Store(key, move, 10, 3, Upper)
	tt.Store(key, move, 20, 6, Lower) // deeper depth - overwrite

	entry, _ := tt.Probe(key)
	if entry.Score != 20 || entry.Depth != 6 || entry.Flag != Lower {
		t.Errorf("expected deeper entry to overwrite, got score=%d depth=%d flag=%v", entry.Score, entry.Depth, entry.Flag)
	}
}

func TestStoreKeeps(t *testing.T) {
	var tt Table
	var move board.Move
	key := uint64(7)

	tt.Store(key, move, 20, 6, Lower)
	tt.Store(key, move, 10, 3, Upper) // shallower depth - don't overwrite

	entry, _ := tt.Probe(key)
	if entry.Score != 20 || entry.Depth != 6 || entry.Flag != Lower {
		t.Errorf("expected original deeper entry preserved, got score=%d depth=%d flag=%v", entry.Score, entry.Depth, entry.Flag)
	}
}

func TestStoreCollision(t *testing.T) {
	var tt Table
	var move board.Move
	keyA := uint64(11)
	keyB := keyA + uint64(Size) // same index, different key

	if tt.index(keyA) != tt.index(keyB) {
		t.Fatalf("keys do not collide")
	}

	tt.Store(keyA, move, 50, 10, Exact) // deep depth
	tt.Store(keyB, move, 1, 1, Upper)   // shallow, but different key - overwrite

	entry, hit := tt.Probe(keyB)
	if !hit || entry.Key != keyB || entry.Score != 1 || entry.Depth != 1 {
		t.Errorf("expected collision to overwrite regardless of depth, got hit=%v key=%d score=%d depth=%d",
			hit, entry.Key, entry.Score, entry.Depth)
	}

	if _, hitA := tt.Probe(keyA); hitA {
		t.Errorf("expected keyA to be replaced by colliding keyB")
	}
}

func TestClear(t *testing.T) {
	var tt Table
	var move board.Move
	key := uint64(123)

	tt.Store(key, move, 77, 8, Exact)
	if _, hit := tt.Probe(key); !hit {
		t.Fatalf("expected hit before clear")
	}

	tt.Clear()
	if _, hit := tt.Probe(key); hit {
		t.Errorf("expected miss after Clear")
	}
}
