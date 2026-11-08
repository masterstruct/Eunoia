//go:build stats

package stats

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

func abs(v int64) uint64 {
	if v < 0 {
		return uint64(-v)
	}
	return uint64(v)
}
