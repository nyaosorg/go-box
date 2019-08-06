package box

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
)

var AnsiCutter = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

func Print(ctx context.Context, nodes []string, out io.Writer) bool {
	b := New()
	b.Height = 0
	value, _, _ := b.Print(ctx, nodes, 0, out)
	return value
}

func (b *box_t) Print(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	selected, columns, nlines := b.PrintNoLastLineFeed(ctx, nodes, offset, out)

	// append last linefeed.
	if nlines > 0 {
		fmt.Fprintln(out)
	}
	b.Cache = nil
	return selected, columns, nlines
}

func (b *box_t) PrintNoLastLineFeed(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	if len(nodes) <= 0 {
		return true, 0, 0
	}

	maxLen := 1
	for _, finfo := range nodes {
		length := runewidth.StringWidth(AnsiCutter.ReplaceAllString(finfo, ""))
		if length > maxLen {
			maxLen = length
		}
	}
	nodePerLine := (b.Width - 1) / (maxLen + 1)
	if nodePerLine <= 0 {
		nodePerLine = 1
	}
	nlines := (len(nodes) + nodePerLine - 1) / nodePerLine

	lines := make([][]byte, nlines)
	row := 0
	for _, finfo := range nodes {
		lines[row] = append(lines[row], finfo...)
		w := runewidth.StringWidth(AnsiCutter.ReplaceAllString(finfo, ""))
		if maxLen < b.Width {
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
	if b.Height > 0 {
		if i_end >= offset+b.Height {
			i_end = offset + b.Height
		}
	}

	if b.Cache == nil {
		b.Cache = make([][]byte, b.Height)
	}
	i := offset
	y := 0
	for {
		if y >= len(b.Cache) {
			b.Cache = append(b.Cache, []byte{})
		}
		// assertion
		if i >= len(lines) {
			panic(fmt.Sprintf("len(lines)==%d i==%d", len(lines), i))
		}
		if bytes.Compare(lines[i], b.Cache[y]) != 0 {
			fmt.Fprint(out, strings.TrimRight(string(lines[i]), " "))
			if len(b.Cache[y]) > 0 {
				fmt.Fprint(out, ERASE_LINE)
			}
			b.Cache[y] = lines[i]
		}
		y++
		if ctx != nil {
			select {
			case <-ctx.Done():
				return false, nodePerLine, nlines
			default:
			}
		}
		i++
		if i >= i_end {
			break
		}
		fmt.Fprintln(out)
	}
	return true, nodePerLine, nlines
}

const (
	CURSOR_OFF = "\x1B[?25l"
	CURSOR_ON  = "\x1B[?25h"
	BOLD_ON    = "\x1B[0;47;30m"
	BOLD_ON2   = "\x1B[0;1;7m"
	BOLD_OFF   = "\x1B[0m"
	UP_N       = "\x1B[%dA"
	ERASE_LINE = "\x1B[0K"
)

func truncate(s string, w int) string {
	return runewidth.Truncate(strings.TrimSpace(s), w, "")
}

const (
	NONE       = 0
	LEFT       = 1
	DOWN       = 2
	UP         = 3
	RIGHT      = 4
	ENTER      = 5
	LEAVE      = 6
	SELECT     = 7
	BACKSELECT = 8
)

type nodeT struct {
	Index int
	Text  string
}

// Choice returns selected string
func Choice(sources []string, out io.Writer) string {
	n := Choose(sources, out)
	if n < 0 {
		return ""
	}
	return sources[n]
}

func ChoiceMulti(sources []string, out io.Writer) []string {
	list := ChooseMulti(sources, out)
	values := make([]string, 0, len(list))
	for _, index := range list {
		values = append(values, sources[index])
	}
	return values
}

// Choice returns the index of selected string
func ChooseMulti(sources []string, out io.Writer) []int {
	cursor := 0
	selected := make(map[int]struct{})

	nodes := make([]*nodeT, 0, len(sources))
	draws := make([]string, 0, len(sources))
	b := New()
	defer b.Close()
	for i, text := range sources {
		val := truncate(text, b.Width-1)
		if val != "" {
			nodes = append(nodes, &nodeT{Index: i, Text: val})
			draws = append(draws, val)
		}
	}
	io.WriteString(out, CURSOR_OFF)
	defer io.WriteString(out, CURSOR_ON)

	if len(nodes) <= 0 {
		nodes = []*nodeT{&nodeT{-1, ""}}
		draws = []string{""}
	}

	offset := 0
	for {
		for index := range selected {
			draws[index] = BOLD_ON + truncate(nodes[index].Text, b.Width-2) + BOLD_OFF
		}
		draws[cursor] = BOLD_ON2 + truncate(nodes[cursor].Text, b.Width-2) + BOLD_OFF
		status, _, h := b.PrintNoLastLineFeed(nil, draws, offset, out)
		if !status {
			return []int{}
		}
		for index, _ := range selected {
			draws[index] = truncate(nodes[index].Text, b.Width-2)
		}
		draws[cursor] = truncate(nodes[cursor].Text, b.Width-2)
		last := cursor
		for last == cursor {
			if bw, ok := out.(*bufio.Writer); ok {
				bw.Flush()
			}
			switch b.GetCmd() {
			case LEFT:
				cursor = (cursor + len(nodes) - h) % len(nodes)
			case RIGHT:
				cursor = (cursor + h) % len(nodes)
			case SELECT:
				if _, ok := selected[cursor]; ok {
					delete(selected, cursor)
				} else {
					selected[cursor] = struct{}{}
				}
				fallthrough
			case DOWN:
				cursor++
				if cursor >= len(nodes) {
					cursor = 0
				}
			case UP:
				cursor--
				if cursor < 0 {
					cursor = len(nodes) - 1
				}
			case BACKSELECT:
				cursor--
				if cursor < 0 {
					cursor = len(nodes) - 1
				}
				if _, ok := selected[cursor]; ok {
					delete(selected, cursor)
				} else {
					selected[cursor] = struct{}{}
				}
			case ENTER:
				selected[cursor] = struct{}{}
				result := make([]int, 0, len(selected))
				for index := range selected {
					result = append(result, index)
				}
				sort.Ints(result)
				return result

			case LEAVE:
				return []int{}
			}

			// x := cursor / h
			y := cursor % h
			if y < offset {
				offset = y
				// offset--
			} else if y >= offset+b.Height {
				offset = y - b.Height + 1
				// offset++
			}
		}
		if h < b.Height {
			if h > 1 {
				fmt.Fprintf(out, UP_N, h-1)
			}
		} else {
			if b.Height > 1 {
				fmt.Fprintf(out, UP_N, b.Height-1)
			}
		}
		fmt.Fprint(out, "\r")
	}
}

func Choose(sources []string, out io.Writer) int {
	selected := ChooseMulti(sources, out)
	if selected == nil || len(selected) <= 0 {
		return -1
	}
	return selected[0]
}
