//go:build !stats

package stats

// here we define methods the program calls
// if it is not compiled with `stats` tag.

type TrackedValue struct{}

func (tv *TrackedValue) Record(value int64) {}

func (tv *TrackedValue) Stats() ValueStats { return ValueStats{} }

type HitMissCounter struct{}

func (c *HitMissCounter) Stats() HitMissStats { return HitMissStats{} }
