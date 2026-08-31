package search

import (
	"time"

	"github.com/masterstruct/Eunoia/internal/tt"
)

type SearchState struct {
	Quiet bool // avoid printing output

	Stop      bool
	Nodes     uint64
	MaxNodes  uint64
	SoftNodes uint64
	StartTime time.Time
	MaxTime   time.Time
	SoftTime  time.Time

	tt *tt.Table
	pv *PVTable

	keyHistory  []uint64 // history of position hashes for 3fold detection
	rootHistLen int

	butterflyHistory *[2][64][64]int
}

func (ss *SearchState) Reset() {
	ss.Quiet = false
	ss.Stop = false
	ss.Nodes = 0
	ss.MaxNodes = 0
	ss.SoftNodes = 0
	ss.StartTime = time.Now()
	ss.MaxTime = time.Time{}
	ss.SoftTime = time.Time{}
}

func (ss *SearchState) ClearTT() {
	ss.tt.Clear()
}

func (ss *SearchState) ClearButterflyHistory() {
	ss.butterflyHistory = &[2][64][64]int{}
}

func (ss *SearchState) Init() {
	ss.tt = &tt.Table{}
	ss.pv = &PVTable{}
	ss.butterflyHistory = &[2][64][64]int{}
}

func (ss *SearchState) SetHistory(gameHistory []uint64) {
	ss.keyHistory = append(ss.keyHistory[:0], gameHistory...)
	ss.rootHistLen = len(ss.keyHistory)
}
