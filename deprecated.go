package box

import (
	"context"
	"io"
)

func New() *BoxT {
	val, err := NewBox()
	if err != nil {
		panic(err.Error())
	}
	return val
}

func Print(ctx context.Context, nodes []string, out io.Writer) bool {
	return PrintX(ctx, nodes, out) == nil
}

func (b *BoxT) Print(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	columns, nlines, err := b.PrintX(ctx, nodes, offset, out)
	return err == nil, columns, nlines
}

func (b *BoxT) PrintNoLastLineFeed(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	if ctx == nil {
		ctx = context.TODO()
	}
	col, row, err := b.PrintNoLastLineFeedX(ctx, nodes, offset, out)
	return err == nil, col, row
}

// Choice returns selected string
func Choice(sources []string, out io.Writer) string {
	val, err := ChoiceX(sources, out)
	if err != nil {
		panic(err.Error())
	}
	return val
}

func ChoiceMulti(sources []string, out io.Writer) []string {
	val, err := ChoiceMultiX(sources, out)
	if err != nil {
		panic(err.Error())
	}
	return val
}

// Choice returns the index of selected string
func ChooseMulti(sources []string, out io.Writer) []int {
	val, err := ChooseMultiX(sources, out)
	if err != nil {
		panic(err.Error())
	}
	return val
}

func Choose(sources []string, out io.Writer) int {
	val, err := ChooseX(sources, out)
	if err != nil {
		panic(err.Error())
	}
	return val
}
