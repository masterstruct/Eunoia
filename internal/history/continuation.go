package history

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

var ContPlies = [...]uint16{1, 2, 4} // mock plies

const plyCount = len(ContPlies)

var contBonuses = [plyCount][5]int32{
	{ // 1
		-250, // bonus base
		300,  // bonus scale
		250,  // malus base
		-300, // malus scale
		2048, // max   bonus
	},
	{ // 2
		-250, // bonus base
		300,  // bonus scale
		250,  // malus base
		-300, // malus scale
		2048, // max   bonus
	},
	{ // 4
		-250, // bonus base
		300,  // bonus scale
		250,  // malus base
		-300, // malus scale
		2048, // max   bonus
	},
}

type ContEntry [2][6][64]int16    // [stm][piece][to]
type ContHist [2][6][64]ContEntry // [prev stm][prev piece][prev to][stm][piece][to]

func (r *ContEntry) get(stm board.Color, pt board.PieceType, to board.Square) int32 {
	if r == nil {
		return 0
	}
	return int32(r[stm][pt][to])
}

func (r *ContEntry) update(stm board.Color, pt board.PieceType, to board.Square, bonus, maxBonus int32) {
	if r != nil {
		gravity(&r[stm][pt][to], bonus, maxBonus, maxHistory)
	}
}

// contains a ContEntry for each ply-N,
// basically ss.stack[ply-N].cont
type ContIndices [plyCount]*ContEntry

// sums the score of all conthist plies
func (i *ContIndices) score(pos *board.Position, move board.Move) (score int32) {
	attacker := pos.Attacker(move)
	stm := pos.SideToMove
	to := move.To()

	for ply := range i {
		score += i[ply].get(stm, attacker, to)
	}
	return
}

// updates all conthist plies at the same time
func (i *ContIndices) update(stm board.Color, pt board.PieceType, to board.Square, depth int32, bonus bool) {
	offset := 2
	if bonus {
		offset = 0
	}
	for ply := range i {
		bonus := contBonuses[ply][0+offset] + contBonuses[ply][1+offset]*depth
		i[ply].update(stm, pt, to, bonus, contBonuses[ply][4])
	}
}
