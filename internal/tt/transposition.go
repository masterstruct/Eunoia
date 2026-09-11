package tt

import (
	"math/bits"
	"unsafe"

	"github.com/masterstruct/Eunoia/internal/board"
)

type Flag uint8

const (
	Exact Flag = iota
	Upper
	Lower
)

type Entry struct {
	Key   uint64
	Move  board.Move
	Score int16
	Depth uint8
	Flag  Flag
}

const DefaultSizeMiB uint = 64

type Table struct {
	entries []Entry
	mask    uint64
}

func NewTable(sizeMiB uint) *Table {
	if sizeMiB == 0 {
		sizeMiB = DefaultSizeMiB
	}
	size := sizeFromMiB(sizeMiB)
	return &Table{
		entries: make([]Entry, size),
		mask:    uint64(size - 1),
	}
}

func (tt *Table) Resize(sizeMiB uint) {
	if sizeMiB == 0 {
		sizeMiB = DefaultSizeMiB
	}
	size := sizeFromMiB(sizeMiB)
	tt.entries = make([]Entry, size)
	tt.mask = uint64(size - 1)
}

func (tt *Table) Clear() {
	if tt != nil {
		clear(tt.entries)
	}
}

func (tt *Table) Store(key uint64, move board.Move, score int16, depth uint8, flag Flag) {
	if tt == nil || len(tt.entries) == 0 {
		return
	}
	entry := &tt.entries[tt.index(key)]
	if entry.Key != key || depth >= entry.Depth {
		entry.Key = key
		entry.Move = move
		entry.Score = score
		entry.Depth = depth
		entry.Flag = flag
	}
}

func (tt *Table) Probe(key uint64) (Entry, bool) {
	if tt == nil || len(tt.entries) == 0 {
		return Entry{}, false
	}
	entry := tt.entries[tt.index(key)]
	return entry, entry.Key == key
}

func (tt *Table) index(key uint64) uint64 {
	return key & tt.mask
}

func sizeFromMiB(mb uint) uint {
	bytes := mb * 1024 * 1024
	entries := bytes / uint(unsafe.Sizeof(Entry{}))
	return nextPow2(entries)
}

func nextPow2(x uint) uint {
	if x <= 1 {
		return 1
	}
	return 1 << bits.Len(x-1)
}
