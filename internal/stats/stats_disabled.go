//go:build !stats

package stats

// here we define methods the program calls
// if it is not compiled with `stats` tag.

type TrackedValue struct{}
