package uci

import (
	"context"
	"fmt"
	"io"
)

func Loop(ctx context.Context, r io.Reader, w io.Writer) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Fprintln(w, "Eunoia chess engine")
	}
}
