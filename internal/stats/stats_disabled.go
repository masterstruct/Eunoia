//go:build !stats

package stats

import (
	"io"
)

// here we define methods the program calls
// if it is not compiled with `stats` tag.

type TrackedValue struct{}

func (tv *TrackedValue) Record(value int64) {}

func (tv *TrackedValue) Stats() ValueStats { return ValueStats{} }

func (tv *TrackedValue) Print() {}

func (s ValueStats) Print() {}

func (s *ValueStats) WriteTo(w io.Writer) (int64, error) { return 0, nil }

type HitMissCounter struct{}

func (c *HitMissCounter) Hit() {}

func (c *HitMissCounter) Miss() {}

func (c *HitMissCounter) Stats() HitMissStats { return HitMissStats{} }
