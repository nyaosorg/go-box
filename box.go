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

var reduntantColorChangePattern = regexp.MustCompile("(\x1B[^m]+m).*?(\x1B[^m]+m)")

func cutReduntantColorChange(s string) string {
	for {
		m := reduntantColorChangePattern.FindStringSubmatchIndex(s)
		if m == nil || len(m) <= 0 {
			return s
		}
		// all = s[m[0]:m[1]]
		first := s[m[2]:m[3]]
		second := s[m[4]:m[5]]

		if first == second {
			s = s[:m[4]] + s[m[5]:]
		} else {
			return s[:m[4]] + cutReduntantColorChange(s[m[4]:])
		}
	}
}

var wtRuneWidth *runewidth.Condition

func init() {
	wtRuneWidth = runewidth.NewCondition()
	if os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != "" {
		wtRuneWidth.EastAsianWidth = false
	}
}

var AnsiCutter = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

func PrintX(ctx context.Context, nodes []string, out io.Writer) error {
	b, err := NewBox()
	if err != nil {
		return err
	}
	b.height = 0
	_, _, err = b.PrintX(ctx, nodes, 0, out)
	return err
}

func (b *BoxT) PrintX(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (int, int, error) {

	columns, nlines, err := b.PrintNoLastLineFeedX(ctx, nodes, offset, out)
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

func (b *BoxT) PrintNoLastLineFeedX(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (int, int, error) {

	if len(nodes) <= 0 {
		return 0, 0, nil
	}

	maxLen := 1
	for _, finfo := range nodes {
		length := wtRuneWidth.StringWidth(AnsiCutter.ReplaceAllString(finfo, ""))
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
		w := wtRuneWidth.StringWidth(AnsiCutter.ReplaceAllString(finfo, ""))
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
				fmt.Fprint(out, _ERASE_LINE)
			}
			b.cache[y] = lines[i]
		}
		y++
		select {
		case <-ctx.Done():
			return nodePerLine, nlines, ctx.Err()
		default:
		}
		i++
		if i >= i_end {
			break
		}
		fmt.Fprintln(out)
	}
	return nodePerLine, nlines, nil
}

const (
	_CURSOR_OFF = "\x1B[?25l"
	_CURSOR_ON  = "\x1B[?25h"
	_BOLD_ON    = "\x1B[0;47;30m"
	_BOLD_ON2   = "\x1B[0;1;7m"
	_BOLD_OFF   = "\x1B[0m"
	_UP_N       = "\x1B[%dA"
	_ERASE_LINE = "\x1B[0K"
)

func truncate(s string, w int) string {
	return wtRuneWidth.Truncate(strings.TrimSpace(s), w, "")
}

type nodeT struct {
	Index int
	Text  string
}

func ChoiceX(sources []string, out io.Writer) (string, error) {
	n, err := ChooseX(sources, out)
	if err != nil {
		return "", err
	}
	if n < 0 {
		return "", nil
	}
	return sources[n], nil
}

func ChoiceMultiX(sources []string, out io.Writer) ([]string, error) {
	list, err := ChooseMultiX(sources, out)
	if err != nil {
		return nil, err
	}
	values := make([]string, 0, len(list))
	for _, index := range list {
		values = append(values, sources[index])
	}
	return values, nil
}

func ChooseMultiX(sources []string, out io.Writer) ([]int, error) {
	cursor := 0
	selected := make(map[int]struct{})

	nodes := make([]*nodeT, 0, len(sources))
	draws := make([]string, 0, len(sources))
	b, err := NewBox()
	if err != nil {
		return nil, err
	}
	defer b.Close()
	for i, text := range sources {
		val := truncate(text, b.width-1)
		if val != "" {
			nodes = append(nodes, &nodeT{Index: i, Text: val})
			draws = append(draws, val)
		}
	}
	io.WriteString(out, _CURSOR_OFF)
	defer io.WriteString(out, _CURSOR_ON)

	if len(nodes) <= 0 {
		nodes = []*nodeT{&nodeT{-1, ""}}
		draws = []string{""}
	}

	ctx := context.TODO()
	offset := 0
	for {
		for index := range selected {
			draws[index] = _BOLD_ON + truncate(nodes[index].Text, b.width-2) + _BOLD_OFF
		}
		draws[cursor] = _BOLD_ON2 + truncate(nodes[cursor].Text, b.width-2) + _BOLD_OFF
		status, _, h := b.PrintNoLastLineFeed(ctx, draws, offset, out)
		if !status {
			return []int{}, nil
		}
		for index := range selected {
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
			case "h", _K_CTRL_B, _K_LEFT, _K_SHIFT_TAB:
				cursor = (cursor + len(nodes) - h) % len(nodes)
			case "H", _K_CTRL_LEFT:
				cursor = (cursor + len(nodes) - h) % len(nodes)
				doSelect()
			case "L", _K_CTRL_RIGHT:
				doSelect()
				fallthrough
			case "l", _K_CTRL_F, _K_RIGHT, "\t":
				cursor = (cursor + h) % len(nodes)
			case " ", "J", _K_CTRL_DOWN:
				doSelect()
				fallthrough
			case "j", _K_CTRL_N, _K_DOWN:
				cursor = (cursor + 1) % len(nodes)
			case "k", _K_CTRL_P, _K_UP:
				cursor = (cursor + len(nodes) - 1) % len(nodes)
			case "\b", "K", _K_CTRL_UP:
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
			case "\x1B", _K_CTRL_G:
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
				fmt.Fprintf(out, _UP_N, h-1)
			}
		} else {
			if b.height > 1 {
				fmt.Fprintf(out, _UP_N, b.height-1)
			}
		}
		fmt.Fprint(out, "\r")
	}
}

func ChooseX(sources []string, out io.Writer) (int, error) {
	selected, err := ChooseMultiX(sources, out)
	if err != nil {
		return -1, err
	}
	if selected == nil || len(selected) <= 0 {
		return -1, nil
	}
	return selected[0], nil
}
