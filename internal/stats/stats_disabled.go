//go:build !stats

package stats

// here we define methods the program calls
// if it is not compiled with `stats` tag.

type TrackedValue struct{}

type ValueStats struct {
	Name     string
	Count    uint64
	Total    int64
	TotalAbs uint64
	TotalSq  uint64
	Min      int64
	Max      int64
	Mean     float64
	StdDev   float64
}

func (tv *TrackedValue) Record(value int64) {}

func (tv *TrackedValue) Stats() ValueStats { return ValueStats{} }
