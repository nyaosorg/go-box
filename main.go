package box

import (
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
	b, err := New()
	if err != nil {
		return nil, err
	}
	r, err := b.SelectIndex(sources, multi, out)
	b.Close()
	return r, err
}

// SelectString returns the strings that user selected.
func SelectString(sources []string, multi bool, out io.Writer) ([]string, error) {
	b, err := New()
	if err != nil {
		return nil, err
	}
	r, err := b.SelectString(sources, multi, out)
	b.Close()
	return r, err
}
