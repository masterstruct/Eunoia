//go:build !stats

package stats // ignore this warning

import (
	"io"
)

// here we define methods the program calls
// if it is not compiled with `stats` tag.

type TrackedValue struct{}

func NewTrackedValue(name string) TrackedValue { return TrackedValue{} }

func (tv *TrackedValue) Record(value int64) {}

func (tv *TrackedValue) Reset() {}

func (tv *TrackedValue) Stats() ValueStats { return ValueStats{} }

func (tv *TrackedValue) Print() {}

func (s ValueStats) Print() {}

func (s *ValueStats) WriteTo(w io.Writer) (int64, error) { return 0, nil }

type HitMissCounter struct{}

func NewHitMissCounter(name string) HitMissCounter { return HitMissCounter{} }

func (c *HitMissCounter) Hit() {}

func (c *HitMissCounter) Miss() {}

func (c *HitMissCounter) Reset() {}

func (c *HitMissCounter) Stats() HitMissStats { return HitMissStats{} }

func (c *HitMissCounter) Print() {}

func (s HitMissStats) Print() {}

func (s *HitMissStats) WriteTo(w io.Writer) (int64, error) { return 0, nil }
