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

func TTSizeFromMB(mb uint) uint {
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
