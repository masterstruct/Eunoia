package main

import (
	"fmt"

	"github.com/masterstruct/Eunoia/internal/board"
)

func init() {
	board.InitBitboards()
}

func main() {
	fmt.Println("Eunoia chess engine")
}
