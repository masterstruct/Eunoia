package search

import (
	"github.com/masterstruct/Eunoia/internal/board"
	"github.com/masterstruct/Eunoia/internal/history"
	"github.com/masterstruct/Eunoia/internal/tt"
)

type SearchState struct {
	Quiet bool // avoid printing output

	TimeManager

	bestMove          board.Move
	bestScore         int32
	bestMoveStability uint8
	depth             int
	selDepth          uint16
	nmpMinPly         uint16

	tt *tt.Table
	pv *PVTable

	Keys  Keys // history of position hashes for 3fold detection
	Stack NodeStack

	history  *history.History
	corrHist *history.CorrHist
}

func (ss *SearchState) Init(ttSizeMiB uint) {
	ss.depth = 1
	ss.tt = tt.NewTable(ttSizeMiB)
	ss.pv = NewPVTable()
	ss.Keys = Keys{keys: make([]uint64, 0, MaxPly)}
	ss.Stack.reset()
	ss.history = &history.History{}
	ss.corrHist = &history.CorrHist{}
}

func (ss *SearchState) PrepareForSearch() {
	ss.Quiet = false
	ss.bestMove = board.NullMove
	ss.bestScore = 0
	ss.bestMoveStability = 0
	ss.depth = 1
	ss.resetTimeManager()
	ss.Stack.reset()
}

func (ss *SearchState) ClearTables() {
	ss.ClearTT()
	*ss.pv = PVTable{}
	ss.history.Clear()
	ss.corrHist.Clear()
}

// clears ONLY the TT - intended for `setoption name Clear Hash`.
// To clear everything use ss.ClearTables()
func (ss *SearchState) ClearTT() {
	ss.tt.Clear()
}

func (ss *SearchState) ResizeTT(sizeMiB uint) {
	ss.tt.Resize(sizeMiB)
}

func (ss *SearchState) PreviousMove(ply uint16, back uint16) (board.Piece, board.Move, bool) {
	if ply >= back {
		n := &ss.Stack[ply-back]
		return n.piece, n.move, n.piece != board.NoPiece
	}
	return board.NoPiece, board.NullMove, false
}

func (ss *SearchState) NewContIndices(ply uint16) history.ContIndices {
	return history.NewContIndices(
		history.NewContEntry(ss.PreviousMove(ply, 1)),
	)
}
