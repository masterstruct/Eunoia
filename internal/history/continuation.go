package history

import (
	"github.com/masterstruct/Eunoia/internal/board"
)

type contParams struct {
	ply                   uint16
	bonusBase, bonusScale int32
	malusBase, malusScale int32
	maxBonus              int32
}

var contPlies = [...]contParams{
	{
		ply:        1,
		bonusBase:  -250,
		bonusScale: 300,
		malusBase:  250,
		malusScale: -300,
		maxBonus:   2048,
	},
	{
		ply:        2,
		bonusBase:  -250,
		bonusScale: 300,
		malusBase:  250,
		malusScale: -300,
		maxBonus:   2048,
	},
}

const plyCount = len(contPlies)

func ContPlyAt(i int) uint16 {
	return contPlies[i].ply
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
	for ply, params := range contPlies {
		base, scale := params.malusBase, params.malusScale
		if bonus {
			base, scale = params.bonusBase, params.bonusScale
		}
		i[ply].update(stm, pt, to, base+scale*depth, params.maxBonus)
	}
}
