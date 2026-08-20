package main

import (
	"os"

	"github.com/masterstruct/Eunoia/internal/uci"
)

func main() {
	uci.Loop(os.Stdin, os.Stdout)
}
