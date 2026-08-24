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

const Size uint = 1 << 22 // 64MiB
const Mask uint = Size - 1

type Table [Size]Entry

func (tt *Table) Clear() {
	*tt = Table{}
}

func (tt *Table) Store(key uint64, move board.Move, score int16, depth uint8, flag Flag) {
	entry := &tt[tt.index(key)]
	if entry.Key != key || depth >= entry.Depth {
		entry.Key = key
		entry.Move = move
		entry.Score = score
		entry.Depth = depth
		entry.Flag = flag
	}
}

func (tt *Table) Probe(key uint64) (Entry, bool) {
	entry := tt[tt.index(key)]
	return entry, entry.Key == key
}

func (t *Table) index(key uint64) uint64 {
	return key & uint64(Mask)
}

func SizeFromMB(mb uint) uint {
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
