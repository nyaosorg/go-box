package box

import (
	"context"
	"io"
)

// Deprecated:
type BoxT = Box

// Deprecated:
func New() *Box {
	val, err := NewBox()
	if err != nil {
		panic(err.Error())
	}
	return val
}

// Deprecated:
func Print(ctx context.Context, nodes []string, out io.Writer) bool {
	return PrintX(ctx, nodes, out) == nil
}

// Deprecated:
func (b *Box) Print(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	columns, nlines, err := b.PrintX(ctx, nodes, offset, out)
	return err == nil, columns, nlines
}

// Deprecated:
func (b *Box) PrintNoLastLineFeed(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	if ctx == nil {
		ctx = context.TODO()
	}
	col, row, err := b.PrintNoLastLineFeedX(ctx, nodes, offset, out)
	return err == nil, col, row
}

// Deprecated: Choice returns returns the string that user selected.
func Choice(sources []string, out io.Writer) string {
	val, err := SelectString(sources, false, out)
	if err != nil {
		panic(err.Error())
	}
	if len(val) <= 0 {
		return ""
	}
	return val[0]
}

// Deprecated: ChoiceMulti returns the strings that user selected.
func ChoiceMulti(sources []string, out io.Writer) []string {
	val, err := SelectString(sources, true, out)
	if err != nil {
		panic(err.Error())
	}
	return val
}

// Deprecated: Choose Multi returns the indices that user selected.
func ChooseMulti(sources []string, out io.Writer) []int {
	val, err := SelectIndex(sources, true, out)
	if err != nil {
		panic(err.Error())
	}
	return val
}

// Deprecated: Choose returns the index that user selected
func Choose(sources []string, out io.Writer) int {
	val, err := SelectIndex(sources, false, out)
	if err != nil {
		panic(err.Error())
	}
	if len(val) <= 0 {
		return -1
	}
	return val[0]
}
