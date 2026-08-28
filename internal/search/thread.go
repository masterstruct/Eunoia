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

func (ss *SearchState) Init() {
	ss.tt = &tt.Table{}
	ss.pv = &PVTable{}
}
