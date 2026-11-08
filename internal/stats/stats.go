//go:build stats

package stats

import (
	"fmt"
	"io"
	"math"
	"os"
)

// to disable warning in VSCode, add
// "gopls": {
//     "buildFlags": ["-tags=stats"]
// },
// to settings.json

// inspired by Viridithas, thank you Cosmo!
// https://github.com/cosmobobak/viridithas/blob/master/src/stats.rs

type TrackedValue struct {
	name     string // name of the tracked value
	count    uint64 // number of samples
	total    int64  // sum of all values
	totalAbs uint64 // sum of all absolute values
	totalSq  uint64 // sum of squared values
	min      int64  // smallest value
	max      int64  // largest value
}

// record a value
func (tv *TrackedValue) Record(value int64) {
	if tv.count == 0 {
		tv.min, tv.max = value, value
	} else {
		tv.min = min(tv.min, value)
		tv.max = max(tv.max, value)
	}
	tv.count++
	tv.total += value
	tv.totalAbs += abs(value)
	tv.totalSq += uint64(value) * uint64(value)
}

func (tv *TrackedValue) Stats() ValueStats {
	return ValueStats{
		Name:     tv.name,
		Count:    tv.count,
		Total:    tv.total,
		TotalAbs: tv.totalAbs,
		TotalSq:  tv.totalSq,
		Min:      tv.min,
		Max:      tv.max,
		Mean:     tv.mean(),
		StdDev:   tv.stdDev(),
	}
}

func (tv *TrackedValue) Print() {
	tv.Stats().Print()
}

func (s ValueStats) Print() {
	s.WriteTo(os.Stdout)
}

func (s *ValueStats) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w,
		"%s: Count %d, Total %d, TotalAbs %d, Min %d, Max %d, Mean %.3f, StdDev %.3f\n",
		s.Name, s.Count, s.Total, s.TotalAbs, s.Min, s.Max, s.Mean, s.StdDev)
	return int64(n), err
}

func (tv *TrackedValue) mean() float64 {
	if tv.count == 0 {
		return 0
	}
	return float64(tv.total) / float64(tv.count)
}

func (tv *TrackedValue) stdDev() float64 {
	if tv.count == 0 {
		return 0
	}
	n := float64(tv.count)
	mean := float64(tv.total) / n
	variance := float64(tv.totalSq)/n - mean*mean

	return math.Sqrt(math.Max(variance, 0))
}

type HitMissCounter struct {
	hits, misses uint64
}

func (c *HitMissCounter) Hit() {
	c.hits++
}

func (c *HitMissCounter) Miss() {
	c.misses++
}

func (c *HitMissCounter) Stats() HitMissStats {
	return HitMissStats{
		Hits:        c.hits,
		Misses:      c.misses,
		Total:       c.hits + c.misses,
		HitPercent:  percent(c.hits, c.hits+c.misses),
		MissPercent: percent(c.misses, c.hits+c.misses),
	}
}

func percent(count, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return 100 * float64(count) / float64(total)
}

func abs(v int64) uint64 {
	if v < 0 {
		return uint64(-v)
	}
	return uint64(v)
}
