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

type HitMissCounter struct{}

type HitMissStats struct {
	Hits        uint64
	Misses      uint64
	Total       uint64
	HitPercent  float64 // 0-100
	MissPercent float64 // 0-100
}

func (c *HitMissCounter) Stats() HitMissStats { return HitMissStats{} }
