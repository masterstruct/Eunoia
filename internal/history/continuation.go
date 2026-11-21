package history

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

var contPlies = [...]int{1, 2, 4} // mock plies

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

const (
	plyCount = len(contPlies)
)

type ContEntry [6][64]int16       // [piece][to]
type ContHist [2][6][64]ContEntry // [stm][prev piece][prev to][piece][to]

func (r *ContEntry) get(pt board.PieceType, to board.Square) int32 {
	if r == nil {
		return 0
	}
	return int32(r[pt][to])
}

func (r *ContEntry) update(pt board.PieceType, to board.Square, bonus, maxBonus int32) {
	if r == nil {
		return
	}
	gravity(&r[pt][to], bonus, maxBonus, maxHistory)
}

// contains a ContEntry for each ply-N,
// basically ss.stack[ply-N].cont
type ContIndices [plyCount]*ContEntry

func NewContIndices(cont1 *ContEntry) ContIndices {
	return ContIndices{cont1}
}

// sums the score of all conthist plies
func (i *ContIndices) score(pos *board.Position, move board.Move) int32 {
	attacker := pos.Attacker(move)
	to := move.To()
	var score int32
	for ply := range plyCount {
		score += i[ply].get(attacker, to)
	}
	return score
}

// updates all conthist plies at the same time
func (i *ContIndices) update(pt board.PieceType, to board.Square, depth int32, bonus bool) {
	offset := 2
	if bonus {
		offset = 0
	}
	for ply := range i {
		bonus := contBonuses[ply][0+offset] + contBonuses[ply][1+offset]*depth
		i[ply].update(pt, to, bonus, contBonuses[ply][4])
	}
}
