package search

const (
	windowSize   = 35
	initialDelta = 50
)

type aspirationWindow struct {
	alpha int16
	beta  int16
	delta int16 // doubles after every widening
}

func (aw *aspirationWindow) widenDown() {
	aw.alpha = max(aw.alpha-aw.delta, -MATE)
	aw.delta *= 2
}

func (aw *aspirationWindow) widenUp() {
	aw.beta = min(aw.beta+aw.delta, MATE)
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
