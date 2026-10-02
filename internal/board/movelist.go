package board

const MaxMoves = 256

type Movelist struct {
	Moves [MaxMoves]ScoredMove
	Len   int
}

func (ml *Movelist) Add(m Move) {
	ml.Moves[ml.Len] = m.ScoredMove(0)
	ml.Len++
}

func (ml *Movelist) AddScoredMove(m ScoredMove) {
	ml.Moves[ml.Len] = m
	ml.Len++
}

func (ml *Movelist) Swap(a, b int) {
	ml.Moves[a], ml.Moves[b] = ml.Moves[b], ml.Moves[a]
}

func (ml *Movelist) Remove(index int) {
	ml.Len--
	ml.Moves[index] = ml.Moves[ml.Len]
}

func (ml *Movelist) IsEmpty() bool {
	return ml.Len == 0
}
