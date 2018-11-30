package box

import (
	"context"
	"io"

	v2 "github.com/zetamatta/go-box/v2"
)

var AnsiCutter = v2.AnsiCutter

func Print(ctx context.Context, nodes []string, out io.Writer) bool {
	return v2.Print(ctx, nodes, out)
}

func Choice(sources []string, out io.Writer) string {
	return v2.Choice(sources, out)
}

func Choose(sources []string, out io.Writer) int {
	return v2.Choose(sources, out)
}
