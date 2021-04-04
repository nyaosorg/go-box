package box

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
)

var wtRuneWidth *runewidth.Condition

func init() {
	wtRuneWidth = runewidth.NewCondition()
	if os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != "" {
		wtRuneWidth.EastAsianWidth = false
	}
}

var ansiCutter = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

func Print(ctx context.Context, nodes []string, out io.Writer) bool {
	b := New()
	b.height = 0
	value, _, _ := b.Print(ctx, nodes, 0, out)
	return value
}

func (b *BoxT) Print(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	selected, columns, nlines := b.PrintNoLastLineFeed(ctx, nodes, offset, out)

	// append last linefeed.
	if nlines > 0 {
		fmt.Fprintln(out)
	}
	b.cache = nil
	return selected, columns, nlines
}

func (b *BoxT) PrintNoLastLineFeed(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	if len(nodes) <= 0 {
		return true, 0, 0
	}

	maxLen := 1
	for _, finfo := range nodes {
		length := wtRuneWidth.StringWidth(ansiCutter.ReplaceAllString(finfo, ""))
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
		w := wtRuneWidth.StringWidth(ansiCutter.ReplaceAllString(finfo, ""))
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
			panic(fmt.Sprintf("len(lines)==%d i==%d", len(lines), i))
		}
		if bytes.Compare(lines[i], b.cache[y]) != 0 {
			fmt.Fprint(out, strings.TrimRight(string(lines[i]), " "))
			if len(b.cache[y]) > 0 {
				fmt.Fprint(out, ERASE_LINE)
			}
			b.cache[y] = lines[i]
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
	return wtRuneWidth.Truncate(strings.TrimSpace(s), w, "")
}

const (
	NONE         = 0
	LEFT         = 1
	DOWN         = 2
	UP           = 3
	RIGHT        = 4
	ENTER        = 5
	LEAVE        = 6
	SELECT_DOWN  = 7
	SELECT_UP    = 8
	SELECT_LEFT  = 9
	SELECT_RIGHT = 10
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
		val := truncate(text, b.width-1)
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
			draws[index] = BOLD_ON + truncate(nodes[index].Text, b.width-2) + BOLD_OFF
		}
		draws[cursor] = BOLD_ON2 + truncate(nodes[cursor].Text, b.width-2) + BOLD_OFF
		status, _, h := b.PrintNoLastLineFeed(nil, draws, offset, out)
		if !status {
			return []int{}
		}
		for index, _ := range selected {
			draws[index] = truncate(nodes[index].Text, b.width-2)
		}
		draws[cursor] = truncate(nodes[cursor].Text, b.width-2)
		last := cursor

		doSelect := func() {
			if _, ok := selected[cursor]; ok {
				delete(selected, cursor)
			} else {
				selected[cursor] = struct{}{}
			}
		}

		for last == cursor {
			if bw, ok := out.(*bufio.Writer); ok {
				bw.Flush()
			}
			key, err := b.getKey()
			if err != nil {
				continue
			}
			switch key {
			case "h", K_CTRL_B, K_LEFT:
				cursor = (cursor + len(nodes) - h) % len(nodes)
			case "H", K_CTRL_LEFT:
				cursor = (cursor + len(nodes) - h) % len(nodes)
				doSelect()
			case "L", K_CTRL_RIGHT:
				doSelect()
				fallthrough
			case "l", K_CTRL_F, K_RIGHT:
				cursor = (cursor + h) % len(nodes)
			case " ", "J", K_CTRL_DOWN:
				doSelect()
				fallthrough
			case "j", K_CTRL_N, K_DOWN:
				cursor = (cursor + 1) % len(nodes)
			case "k", K_CTRL_P, K_UP:
				cursor = (cursor + len(nodes) - 1) % len(nodes)
			case "\b", "K", K_CTRL_UP:
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
				return result
			case "\x1B", K_CTRL_G:
				return []int{}
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
				fmt.Fprintf(out, UP_N, h-1)
			}
		} else {
			if b.height > 1 {
				fmt.Fprintf(out, UP_N, b.height-1)
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
