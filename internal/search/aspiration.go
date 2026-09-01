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
	aw.alpha -= aw.delta
	aw.delta *= 2
}

func (aw *aspirationWindow) widenUp() {
	aw.beta += aw.delta
	aw.delta *= 2
}

func (aw *aspirationWindow) centerAround(score int16) {
	aw.alpha = score - windowSize
	aw.beta = score + windowSize
	aw.delta = initialDelta
}

func newAspirationWindow() aspirationWindow {
	return aspirationWindow{-INF, +INF, initialDelta}
}
