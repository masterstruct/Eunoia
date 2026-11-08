//go:build stats

// to run tests in VSCode, add
// "go.testFlags": [
//     "-tags=stats"
// ],
// to settings.json

package stats

import (
	"math"
	"testing"
)

func assertValues(t *testing.T, tv *TrackedValue, count uint64, total int64,
	totalAbs, totalSq uint64, min, max int64) {
	t.Helper()
	if tv.count != count {
		t.Errorf("count mismatch: want %d, got %d", count, tv.count)
	}
	if tv.total != total {
		t.Errorf("total mismatch: want %d, got %d", total, tv.total)
	}
	if tv.totalAbs != totalAbs {
		t.Errorf("totalAbs mismatch: want %d, got %d", totalAbs, tv.totalAbs)
	}
	if tv.totalSq != totalSq {
		t.Errorf("totalSq mismatch: want %d, got %d", totalSq, tv.totalSq)
	}
	if tv.min != min {
		t.Errorf("min mismatch: want %d, got %d", min, tv.min)
	}
	if tv.max != max {
		t.Errorf("max mismatch: want %d, got %d", max, tv.max)
	}
}

func TestRecord(t *testing.T) {
	tests := []struct {
		name              string
		values            []int64
		count             uint64
		total             int64
		totalAbs, totalSq uint64
		min, max          int64
	}{
		{
			name:   "single positive value",
			values: []int64{64},
			count:  1, total: 64, totalAbs: 64, totalSq: 4096,
			min: 64, max: 64,
		},
		{
			name:   "single negative value",
			values: []int64{-3},
			count:  1, total: -3, totalAbs: 3, totalSq: 9,
			min: -3, max: -3,
		},
		{
			name:   "multiple values",
			values: []int64{-2, 3, 0},
			count:  3, total: 1, totalAbs: 5, totalSq: 13,
			min: -2, max: 3,
		},
		{
			name:   "multiple large values",
			values: []int64{-345978, -2348056, 45678, 507843, 0},
			count:  5, total: -2140513, totalAbs: 3247555, totalSq: 5893058747953,
			min: -2348056, max: 507843,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tv TrackedValue
			for _, value := range tt.values {
				tv.Record(value)
			}
			assertValues(t, &tv, tt.count, tt.total,
				tt.totalAbs, tt.totalSq, tt.min, tt.max)
		})
	}
}

func Test_abs(t *testing.T) {
	tests := []struct {
		name string
		v    int64
		want uint64
	}{
		{
			name: "zero",
			v:    0,
			want: 0,
		},
		{
			name: "negative value becomes positive",
			v:    -234568,
			want: 234568,
		},
		{
			name: "positive value stays unchanged",
			v:    934856,
			want: 934856,
		},
		{
			name: "smallest int64 value",
			v:    math.MinInt64,
			want: math.MaxInt64 + 1,
		},
		{
			name: "largest int64 value",
			v:    math.MaxInt64,
			want: math.MaxInt64 + 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := abs(tt.v)
			if got != tt.want {
				t.Errorf("abs() = %v, want %v", got, tt.want)
			}
		})
	}
}
