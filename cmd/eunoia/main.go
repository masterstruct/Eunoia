package main

import (
	"os"

	"github.com/masterstruct/Eunoia/internal/uci"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "bench" {
		uci.Bench()
		return
	}
	uci.Loop(os.Stdin, os.Stdout)
}
