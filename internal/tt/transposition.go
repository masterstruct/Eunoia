package tt

import (
	"math/bits"
	"runtime/debug"
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

	// for hashfull
	totalEntries uint64
	usedEntries  uint64
}

func NewTable(sizeMiB uint) *Table {
	if sizeMiB == 0 {
		sizeMiB = DefaultSizeMiB
	}
	size := sizeFromMiB(sizeMiB)

	return &Table{
		entries:      make([]Entry, size),
		mask:         size - 1,
		totalEntries: size,
	}
}

func (tt *Table) Resize(sizeMiB uint) {
	if sizeMiB == 0 {
		sizeMiB = DefaultSizeMiB
	}
	size := sizeFromMiB(sizeMiB)

	if size == tt.totalEntries {
		// TT is already correct size - clear
		// it instead of allocating a new one
		tt.Clear()
		return
	}

	tt.entries = make([]Entry, size)
	tt.mask = size - 1
	tt.totalEntries = size
	tt.usedEntries = 0

	// garbage collect the old TT.entries slice
	debug.FreeOSMemory()
}

func (tt *Table) Clear() {
	if tt == nil || tt.usedEntries == 0 {
		return
	}
	clear(tt.entries)
	tt.usedEntries = 0
}

func (tt *Table) Store(key uint64, move board.Move, score int16, depth uint8, flag Flag) {
	if tt == nil || tt.totalEntries == 0 {
		return
	}
	entry := &tt.entries[tt.index(key)]
	if entry.Key != key || depth >= entry.Depth {
		if entry.Key == 0 {
			tt.usedEntries++
		}
		entry.Key = key
		entry.Move = move
		entry.Score = score
		entry.Depth = depth
		entry.Flag = flag
	}
}

func (tt *Table) Probe(key uint64) (Entry, bool) {
	if tt == nil || tt.usedEntries == 0 || tt.totalEntries == 0 {
		return Entry{}, false
	}
	entry := tt.entries[tt.index(key)]
	return entry, entry.Key == key
}

func (tt *Table) index(key uint64) uint64 {
	return key & tt.mask
}

func sizeFromMiB(mb uint) uint64 {
	bytes := mb * 1024 * 1024
	entries := bytes / uint(unsafe.Sizeof(Entry{}))
	return nextPow2(entries)
}

func nextPow2(x uint) uint64 {
	if x <= 1 {
		return 1
	}
	return 1 << bits.Len(x-1)
}

func (tt *Table) Hashfull() uint64 {
	if tt == nil || tt.totalEntries == 0 {
		return 0
	}
	return (tt.usedEntries * 1000) / tt.totalEntries
}
