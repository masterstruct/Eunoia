//go:build stats

package stats

// to disable warning in VSCode, add
// "gopls": {
//     "buildFlags": ["-tags=stats"]
// },
// to settings.json

import "fmt"

func foo() {
	fmt.Println("boo")
}
