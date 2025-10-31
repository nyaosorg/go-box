package box

import (
	"context"
	"io"
)

// Println outputs the given items in a tabular layout and appends
// a newline after the final line.
func Println(nodes []string, out io.Writer) error {
	b, err := New()
	if err != nil {
		return err
	}
	b.Height = 0
	_, _, err = b.Println(nodes, 0, out)
	b.Close()
	return err
}

type nodeT struct {
	Index int
	Text  string
}

// SelectIndex returns the indexes that user selected.
func SelectIndex(sources []string, multi bool, out io.Writer) ([]int, error) {
	return SelectIndexContext(context.TODO(), sources, multi, out)
}

func SelectIndexContext(ctx context.Context, sources []string, multi bool, out io.Writer) ([]int, error) {
	b, err := New()
	if err != nil {
		return nil, err
	}
	r, err := b.SelectIndexContext(ctx, sources, multi, out)
	b.Close()
	return r, err
}

// SelectString returns the strings that user selected.
func SelectString(sources []string, multi bool, out io.Writer) ([]string, error) {
	return SelectStringContext(context.TODO(), sources, multi, out)
}

func SelectStringContext(ctx context.Context, sources []string, multi bool, out io.Writer) ([]string, error) {
	b, err := New()
	if err != nil {
		return nil, err
	}
	r, err := b.SelectStringContext(ctx, sources, multi, out)
	b.Close()
	return r, err
}
