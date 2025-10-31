package box

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/mattn/go-tty"

	"github.com/nyaosorg/go-box/v3/internal/ansi"
	"github.com/nyaosorg/go-box/v3/internal/keys"
)

type Box struct {
	width  int
	height int
	cache  [][]byte
	tty    _Tty
}

func New() (*Box, error) {
	tty1, err := tty.Open()
	if err != nil {
		return nil, err
	}
	w, h, err := tty1.Size()
	return &Box{
		width:  w,
		height: h,
		tty:    _GoTty{TTY: tty1},
	}, err
}

func (b *Box) Close() error {
	return b.tty.Close()
}

// Println outputs the given items in a tabular layout and appends
// a newline after the final line.
func (b *Box) Println(
	nodes []string,
	offset int,
	out io.Writer) (int, int, error) {

	columns, nlines, err := b.Print(nodes, offset, out)
	if err != nil {
		return columns, nlines, err
	}
	// append last linefeed.
	if nlines > 0 {
		fmt.Fprintln(out)
	}
	b.cache = nil
	return columns, nlines, nil
}

// Print outputs the given items in a tabular layout without appending
// a trailing newline. It returns the number of rows and columns printed.
func (b *Box) Print(
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
	nodePerLine := (b.width - 1) / (maxLen + 1)
	if nodePerLine <= 0 {
		nodePerLine = 1
	}
	nlines := (len(nodes) + nodePerLine - 1) / nodePerLine

	lines := make([][]byte, nlines)
	row := 0
	for _, finfo := range nodes {
		lines[row] = append(lines[row], finfo...)
		w := rw.StringWidth(ansi.RxSequence.ReplaceAllString(finfo, ""))
		if maxLen < b.width {
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
	if b.height > 0 {
		if i_end >= offset+b.height {
			i_end = offset + b.height
		}
	}

	if b.cache == nil {
		b.cache = make([][]byte, b.height)
	}
	i := offset
	y := 0
	for {
		if y >= len(b.cache) {
			b.cache = append(b.cache, []byte{})
		}
		// assertion
		if i >= len(lines) {
			return 0, 0, fmt.Errorf("assertion failed: len(lines)==%d i==%d", len(lines), i)
		}
		if !bytes.Equal(lines[i], b.cache[y]) {
			line := strings.TrimRight(string(lines[i]), " ")
			line = cutReduntantColorChange(line)
			io.WriteString(out, line)
			if len(b.cache[y]) > 0 {
				fmt.Fprint(out, ansi.EraseLine)
			}
			b.cache[y] = lines[i]
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

// SelectIndex returns the indexes that user selected.
func (b *Box) SelectIndex(sources []string, multi bool, out io.Writer) ([]int, error) {
	return b.SelectIndexContext(context.TODO(), sources, multi, out)
}

func (b *Box) SelectIndexContext(ctx context.Context, sources []string, multi bool, out io.Writer) ([]int, error) {
	cursor := 0
	selected := make(map[int]struct{})

	nodes := make([]nodeT, 0, len(sources))
	draws := make([]string, 0, len(sources))
	for i, text := range sources {
		val := truncate(text, b.width-1)
		if val != "" {
			nodes = append(nodes, nodeT{Index: i, Text: val})
			draws = append(draws, val)
		}
	}
	io.WriteString(out, ansi.CursorOff)
	defer io.WriteString(out, ansi.CursorOn)

	if len(nodes) <= 0 {
		nodes = []nodeT{nodeT{-1, ""}}
		draws = []string{""}
	}

	offset := 0
	for {
		for index := range selected {
			draws[index] = ansi.BoldOn + truncate(nodes[index].Text, b.width-2) + ansi.BoldOff
		}
		draws[cursor] = ansi.BoldOn2 + truncate(nodes[cursor].Text, b.width-2) + ansi.BoldOff
		_, h, err := b.Print(draws, offset, out)
		if err != nil {
			return []int{}, err
		}
		for index := range selected {
			draws[index] = truncate(nodes[index].Text, b.width-2)
		}
		draws[cursor] = truncate(nodes[cursor].Text, b.width-2)
		last := cursor

		var doSelect func()
		if multi {
			doSelect = func() {
				if _, ok := selected[cursor]; ok {
					delete(selected, cursor)
				} else {
					selected[cursor] = struct{}{}
				}
			}
		} else {
			doSelect = func() {}
		}

		for last == cursor {
			if bw, ok := out.(*bufio.Writer); ok {
				bw.Flush()
			}
			key, err := b.tty.GetKey()
			if err != nil {
				continue
			}
			switch key {
			case "h", keys.CtrlB, keys.Left, keys.ShiftTab:
				cursor = (cursor + len(nodes) - h) % len(nodes)
			case "H", keys.CtrlLeft:
				cursor = (cursor + len(nodes) - h) % len(nodes)
				doSelect()
			case "L", keys.CtrlRight:
				doSelect()
				fallthrough
			case "l", keys.CtrlF, keys.Right, "\t":
				cursor = (cursor + h) % len(nodes)
			case " ", "J", keys.CtrlDown:
				doSelect()
				fallthrough
			case "j", keys.CtrlN, keys.Down:
				cursor = (cursor + 1) % len(nodes)
			case "k", keys.CtrlP, keys.Up:
				cursor = (cursor + len(nodes) - 1) % len(nodes)
			case "\b", "K", keys.CtrlUp:
				cursor = (cursor + len(nodes) - 1) % len(nodes)
				doSelect()
			case "\r", "\n":
				var result []int
				if len(selected) > 0 {
					result = make([]int, 0, len(selected))
					for index := range selected {
						result = append(result, index)
					}
					sort.Ints(result)
				} else {
					result = []int{cursor}
				}
				return result, nil
			case "\x1B", keys.CtrlG:
				return []int{}, nil
			}

			// x := cursor / h
			y := cursor % h
			if y < offset {
				offset = y
				// offset--
			} else if y >= offset+b.height {
				offset = y - b.height + 1
				// offset++
			}
		}
		if h < b.height {
			if h > 1 {
				fmt.Fprintf(out, ansi.UpN, h-1)
			}
		} else {
			if b.height > 1 {
				fmt.Fprintf(out, ansi.UpN, b.height-1)
			}
		}
		fmt.Fprint(out, "\r")
	}
}

// SelectString returns the strings that user selected.
func (b *Box) SelectString(sources []string, multi bool, out io.Writer) ([]string, error) {
	return b.SelectStringContext(context.TODO(), sources, multi, out)

}
func (b *Box) SelectStringContext(ctx context.Context, sources []string, multi bool, out io.Writer) ([]string, error) {
	list, err := b.SelectIndexContext(ctx, sources, multi, out)
	if err != nil {
		return nil, err
	}
	values := make([]string, 0, len(list))
	for _, index := range list {
		values = append(values, sources[index])
	}
	return values, nil
}
