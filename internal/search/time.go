package search

import (
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
)

type TimeManager struct {
	Stop     bool
	MaxDepth int

	// set once in ss.Init and changed with `UpdateMoveOverhead`.
	// should persist between `ucinewgame` calls - do NOT clear
	MoveOverhead uint

	Nodes, MaxNodes, SoftNodes   uint64
	StartTime, MaxTime, SoftTime time.Time
}

type LimitType int

const (
	Soft LimitType = iota
	Hard

	MinMoveOverhead     = 0
	DefaultMoveOverhead = 20
	MaxMoveOverhead     = 5000
)

func (tm *TimeManager) hardLimitReached() bool {
	if tm.Stop {
		return true
	}
	if tm.MaxNodes > 0 && tm.Nodes >= tm.MaxNodes {
		return true
	}
	if tm.Nodes&2047 == 0 && // check hard time limit every 2048 nodes
		!tm.MaxTime.IsZero() && time.Now().After(tm.MaxTime) {
		return true
	}
	return false
}

func (tm *TimeManager) softLimitReached() bool {
	return tm.Stop ||
		(tm.SoftNodes > 0 && tm.Nodes >= tm.SoftNodes) ||
		(!tm.SoftTime.IsZero() && time.Now().After(tm.SoftTime))
}

func (tm *TimeManager) ShouldStop(limitType LimitType) bool {
	if tm.Stop {
		return true
	}

	stop := false

	switch limitType {
	case Soft:
		stop = tm.softLimitReached()
	case Hard:
		stop = tm.hardLimitReached()
	}

	if stop {
		tm.Stop = true
	}
	return stop
}

type GoLimits struct {
	WTime, BTime int64
	WInc, BInc   int64
	MoveTime     int64
	Depth        int64
	Nodes        int64
	Infinite     bool
}

func (tm *TimeManager) resetTimeManager() {
	tm.Stop = false
	tm.MaxDepth = MaxPly
	tm.Nodes = 0
	tm.MaxNodes = 0
	tm.SoftNodes = 0
	tm.StartTime = time.Now()
	tm.MaxTime = time.Time{}
	tm.SoftTime = time.Time{}
}

func (tm *TimeManager) SetLimits(limits GoLimits, sideToMove board.Color) {
	tm.resetTimeManager()

	if limits.Depth <= 0 || limits.Depth > MaxPly || limits.Infinite {
		limits.Depth = MaxPly
	}
	tm.MaxDepth = int(limits.Depth)

	if limits.Nodes > 0 {
		tm.MaxNodes = uint64(limits.Nodes)
	}

	var remainingTime int64
	var increment int64
	if sideToMove == board.Black {
		remainingTime = limits.BTime
		increment = limits.BInc
	} else {
		remainingTime = limits.WTime
		increment = limits.WInc
	}

	if remainingTime <= 0 && limits.MoveTime <= 0 {
		return
	}

	// it is possible to send multiple time constraints, example:
	// go wtime 60000 btime 60000 movetime 5000
	// the engine should stop whenever any limit is reached:
	// 5 seconds (movetime) or internal soft/hard time limits.

	remainingTime = max(remainingTime, 0)
	increment = max(increment, 0)

	var hardTime, softTime int64

	if remainingTime > 0 {
		softTime = max((remainingTime/30+increment*7/10)-int64(tm.MoveOverhead), 1)
		tm.SoftTime = tm.StartTime.Add(time.Duration(softTime) * time.Millisecond)

		hardTime = max(remainingTime/3+increment*7/10, 1)
	}

	if limits.MoveTime > 0 {
		if hardTime == 0 {
			hardTime = limits.MoveTime
		} else {
			hardTime = min(hardTime, limits.MoveTime)
		}
	}

	hardTime = max(hardTime-int64(tm.MoveOverhead), 1)
	tm.MaxTime = tm.StartTime.Add(time.Duration(hardTime) * time.Millisecond)
}

func (tm *TimeManager) UpdateMoveOverhead(ms int) {
	if ms < MinMoveOverhead {
		ms = MinMoveOverhead
	} else if ms > MaxMoveOverhead {
		ms = MaxMoveOverhead
	}

	tm.MoveOverhead = uint(ms)
}
