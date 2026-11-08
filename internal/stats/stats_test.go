//go:build stats

// to run tests in VSCode, add
// "go.testFlags": [
//     "-tags=stats"
// ],
// to settings.json

package stats

import "testing"

func TestFoo(t *testing.T) {
	foo()
}
