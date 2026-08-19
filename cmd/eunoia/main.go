package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/masterstruct/Eunoia/internal/uci"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	uci.Loop(ctx, os.Stdin, os.Stdout)
}
