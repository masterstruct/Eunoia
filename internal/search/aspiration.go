package search

const windowSize = 35

type aspirationWindow struct {
	alpha int16
	beta  int16
}

func (aw *aspirationWindow) widenDown() {
	aw.alpha -= windowSize
}

func (aw *aspirationWindow) widenUp() {
	aw.beta += windowSize
}

func (aw *aspirationWindow) centerAround(score int16) {
	aw.alpha = score - windowSize
	aw.beta = score + windowSize
}

func newAspirationWindow() aspirationWindow {
	return aspirationWindow{-INF, +INF}
}
