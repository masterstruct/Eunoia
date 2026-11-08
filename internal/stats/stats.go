//go:build stats

package stats

// to disable the warning in VSCode
// and properly run tests, add
// "gopls": {
//     "buildFlags": ["-tags=stats"]
// },
// to settings.json

import (
	"fmt"
	"io"
	"math"
	"os"
)

// inspired by Viridithas, thank you Cosmo!
// https://github.com/cosmobobak/viridithas/blob/master/src/stats.rs

type TrackedValue struct {
	Name     string // name of the tracked value
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

func (tv *TrackedValue) Reset() {
	tv.count = 0
	tv.total = 0
	tv.totalAbs = 0
	tv.totalSq = 0
	tv.min = 0
	tv.max = 0
}

func (tv *TrackedValue) Stats() ValueStats {
	return ValueStats{
		Name:     tv.Name,
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
	Name         string
	hits, misses uint64
}

func (c *HitMissCounter) Hit() {
	c.hits++
}

func (c *HitMissCounter) Miss() {
	c.misses++
}

func (c *HitMissCounter) Reset() {
	c.hits = 0
	c.misses = 0
}

func (c *HitMissCounter) Stats() HitMissStats {
	return HitMissStats{
		Name:        c.Name,
		Hits:        c.hits,
		Misses:      c.misses,
		Total:       c.hits + c.misses,
		HitPercent:  percent(c.hits, c.hits+c.misses),
		MissPercent: percent(c.misses, c.hits+c.misses),
	}
}

func (c *HitMissCounter) Print() {
	c.Stats().Print()
}

func (s HitMissStats) Print() {
	s.WriteTo(os.Stdout)
}

func (s *HitMissStats) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w,
		"%s: Total %d, Hits %d (%.2f%%), Misses %d (%.2f%%)\n",
		s.Name, s.Total, s.Hits, s.HitPercent, s.Misses, s.MissPercent)
	return int64(n), err
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
