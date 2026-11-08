package stats

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

type HitMissStats struct {
	Hits        uint64
	Misses      uint64
	Total       uint64
	HitPercent  float64 // 0-100
	MissPercent float64 // 0-100
}
