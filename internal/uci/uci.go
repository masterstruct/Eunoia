package uci

import (
	"context"
	"fmt"
)

func Loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Println("Eunoia chess engine")
	}
}
