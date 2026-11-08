//go:build !stats

package stats

// here we define methods the program calls
// if it not compiled with `stats` tag.

func foo() {}
