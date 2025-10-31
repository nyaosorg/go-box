package grid

import (
	"bytes"
	"fmt"
	"io"

	"strings"

	"github.com/nyaosorg/go-box/v3/internal/ansi"
)

type Grid struct {
	Width  int
	Height int
	cache  [][]byte
}

// Println outputs the given items in a tabular layout and appends
// a newline after the final line.
func (g *Grid) Println(
	nodes []string,
	offset int,
	out io.Writer) (int, int, error) {

	columns, nlines, err := g.Print(nodes, offset, out)
	if err != nil {
		return columns, nlines, err
	}
	// append last linefeed.
	if nlines > 0 {
		fmt.Fprintln(out)
	}
	g.cache = nil
	return columns, nlines, nil
}

// Print outputs the given items in a tabular layout without appending
// a trailing newline. It returns the number of rows and columns printed.
func (g *Grid) Print(
	nodes []string,
	offset int,
	out io.Writer) (int, int, error) {

	if len(nodes) <= 0 {
		return 0, 0, nil
	}

	rw := wtRuneWidth.Value()
	maxLen := 1
	for _, finfo := range nodes {
		length := rw.StringWidth(ansi.RxSequence.ReplaceAllString(finfo, ""))
		if length > maxLen {
			maxLen = length
		}
	}
	nodePerLine := (g.Width - 1) / (maxLen + 1)
	if nodePerLine <= 0 {
		nodePerLine = 1
	}
	nlines := (len(nodes) + nodePerLine - 1) / nodePerLine

	lines := make([][]byte, nlines)
	row := 0
	for _, finfo := range nodes {
		lines[row] = append(lines[row], finfo...)
		w := rw.StringWidth(ansi.RxSequence.ReplaceAllString(finfo, ""))
		if maxLen < g.Width {
			for i := maxLen + 1; i > w; i-- {
				lines[row] = append(lines[row], ' ')
			}
		}
		row++
		if row >= nlines {
			row = 0
		}
	}
	i_end := len(lines)
	if g.Height > 0 {
		if i_end >= offset+g.Height {
			i_end = offset + g.Height
		}
	}

	if g.cache == nil {
		g.cache = make([][]byte, g.Height)
	}
	i := offset
	y := 0
	for {
		if y >= len(g.cache) {
			g.cache = append(g.cache, []byte{})
		}
		// assertion
		if i >= len(lines) {
			return 0, 0, fmt.Errorf("assertion failed: len(lines)==%d i==%d", len(lines), i)
		}
		if !bytes.Equal(lines[i], g.cache[y]) {
			line := strings.TrimRight(string(lines[i]), " ")
			line = cutReduntantColorChange(line)
			io.WriteString(out, line)
			if len(g.cache[y]) > 0 {
				fmt.Fprint(out, ansi.EraseLine)
			}
			g.cache[y] = lines[i]
		}
		y++
		i++
		if i >= i_end {
			break
		}
		fmt.Fprintln(out)
	}
	return nodePerLine, nlines, nil
}
