package search

import (
	"time"

	"github.com/masterstruct/Eunoia/internal/board"
)

type TimeManager struct {
	Stop     bool
	MaxDepth uint16

	// set once in ss.Init and changed with `UpdateMoveOverhead`.
	// should persist between `ucinewgame` calls - do NOT clear
	MoveOverhead uint

	Nodes, MaxNodes, SoftNodes uint64
	StartTime                  time.Time
	MaxTime, SoftTime          time.Duration
}

type LimitType int

const (
	Soft LimitType = iota
	Hard

	MinMoveOverhead     = 0
	DefaultMoveOverhead = 20
	MaxMoveOverhead     = 5000

	softScale    = 33
	softIncScale = 700

	hardScale    = 333
	hardIncScale = 700

	moveStabilityHead     = 1.8
	moveStabilityDecay    = 1.5
	moveStabilityMin      = 0.6
	moveStabilityMinDepth = 5

	scoreStabilityThreshold = 14
	scoreStabilityBase      = 1.23
	scoreStabilityScale     = 0.045
	scoreStabilityMin       = 0.9
)

func (ss *SearchState) hardLimitReached() bool {
	if ss.Stop {
		return true
	}
	if ss.MaxNodes > 0 && ss.Nodes >= ss.MaxNodes {
		return true
	}
	if ss.Nodes&2047 == 0 && // check hard time limit every 2048 nodes
		ss.MaxTime > 0 && time.Since(ss.StartTime) >= ss.MaxTime {
		return true
	}
	if ss.depth > int(ss.MaxDepth) {
		return true
	}
	return false
}

func (ss *SearchState) softLimitReached() bool {
	return ss.Stop ||
		(ss.SoftNodes > 0 && ss.Nodes >= ss.SoftNodes) ||
		ss.depth > int(ss.MaxDepth) ||
		(ss.SoftTime > 0 && time.Since(ss.StartTime) >= ss.scaledSoftTime())
}

func (ss *SearchState) ShouldStop(limitType LimitType) bool {
	// always complete first depth
	if ss.depth <= 1 {
		return false
	}

	if ss.Stop {
		return true
	}

	stop := false

	switch limitType {
	case Soft:
		stop = ss.softLimitReached()
	case Hard:
		stop = ss.hardLimitReached()
	}

	if stop {
		ss.Stop = true
	}
	return stop
}

func (ss *SearchState) updateTM(previousMove board.Move, avgScore, score int32) {
	if previousMove == ss.bestMove {
		ss.bestMoveStability++
	} else {
		ss.bestMoveStability = 0
	}

	if score > avgScore-scoreStabilityThreshold && score < avgScore+scoreStabilityThreshold {
		ss.scoreStability++
	} else {
		ss.scoreStability = 0
	}
}

func (ss *SearchState) scaledSoftTime() time.Duration {
	softTime := ss.SoftTime
	if ss.depth >= moveStabilityMinDepth {
		softTime = time.Duration(float64(softTime) * moveStability(ss.bestMoveStability) * scoreStability(ss.scoreStability))
	}
	if ss.MaxTime > 0 {
		softTime = min(softTime, ss.MaxTime)
	}
	return softTime
}

func moveStability(stability uint8) float64 {
	// not asspulled, tested in desmos
	// c + a / (1 + bx)
	// a - head; b - decay; c - min; x - stability
	// stability starts at 0
	return moveStabilityMin + (moveStabilityHead / (1 + moveStabilityDecay*float64(stability)))
}

func scoreStability(stability uint8) float64 {
	return max(scoreStabilityMin, scoreStabilityBase-float64(stability)*scoreStabilityScale)
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
	tm.MaxTime = 0
	tm.SoftTime = 0
}

func (tm *TimeManager) SetLimits(limits GoLimits, stm board.Color) {
	tm.resetTimeManager()

	if limits.Depth <= 0 || limits.Depth > int64(MaxPly) || limits.Infinite {
		limits.Depth = int64(MaxPly)
	}
	tm.MaxDepth = uint16(limits.Depth)

	if limits.Nodes > 0 {
		tm.MaxNodes = uint64(limits.Nodes)
	}

	var remainingTime int64
	var increment int64
	if stm == board.Black {
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
		softTime = remainingTime*softScale/1000 + increment*softIncScale/1000
		softTime = max(softTime-int64(tm.MoveOverhead), 1)
		tm.SoftTime = time.Duration(softTime) * time.Millisecond

		hardTime = remainingTime*hardScale/1000 + increment*hardIncScale/1000
		hardTime = max(hardTime, 1)
	}

	if limits.MoveTime > 0 {
		if hardTime == 0 {
			hardTime = limits.MoveTime
		} else {
			hardTime = min(hardTime, limits.MoveTime)
		}
	}

	if remainingTime > 0 {
		hardTime = min(hardTime, max(remainingTime-int64(tm.MoveOverhead), 1))
	}

	tm.MaxTime = time.Duration(hardTime) * time.Millisecond
}

func (tm *TimeManager) UpdateMoveOverhead(ms int) {
	if ms < MinMoveOverhead {
		ms = MinMoveOverhead
	} else if ms > MaxMoveOverhead {
		ms = MaxMoveOverhead
	}

	tm.MoveOverhead = uint(ms)
}
