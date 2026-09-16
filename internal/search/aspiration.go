package search

const (
	aspirationMinDepth = 6
	windowSize         = 35
	initialDelta       = 50
)

type aspirationWindow struct {
	alpha int16
	beta  int16
	delta int16 // doubles after every widening
}

func (aw *aspirationWindow) widenDown(score int16) {
	aw.alpha = max(score-aw.delta, -MATE)
	aw.delta *= 2
}

func (aw *aspirationWindow) widenUp(score int16) {
	aw.beta = min(score+aw.delta, MATE)
	aw.delta *= 2
}

func (aw *aspirationWindow) centerAround(score int16) {
	aw.alpha = max(score-windowSize, -MATE)
	aw.beta = min(score+windowSize, MATE)
	aw.delta = initialDelta
}

func newAspirationWindow() aspirationWindow {
	return aspirationWindow{-INF, +INF, initialDelta}
}
