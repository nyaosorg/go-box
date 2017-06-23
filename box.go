package box

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

var ansiCutter = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

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
	fmt.Fprintln(out)
	b.Cache = nil
	return selected, columns, nlines
}

func (b *box_t) PrintNoLastLineFeed(ctx context.Context,
	nodes []string,
	offset int,
	out io.Writer) (bool, int, int) {

	maxLen := 1
	for _, finfo := range nodes {
		length := runewidth.StringWidth(ansiCutter.ReplaceAllString(finfo, ""))
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
		w := runewidth.StringWidth(ansiCutter.ReplaceAllString(finfo, ""))
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
		if bytes.Compare(lines[i], b.Cache[y]) != 0 {
			b.Cache[y] = lines[i]
			fmt.Fprint(out, string(lines[i]))
			fmt.Fprint(out, "\x1B[0K")
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
	BOLD_OFF   = "\x1B[0m"

	K_LEFT  = 0x25
	K_RIGHT = 0x27
	K_UP    = 0x26
	K_DOWN  = 0x28
)

func truncate(s string, w int) string {
	return runewidth.Truncate(strings.TrimSpace(s), w, "")
}

const (
	NONE  = 0
	LEFT  = 1
	DOWN  = 2
	UP    = 3
	RIGHT = 4
	ENTER = 5
	LEAVE = 6
)

func Choice(sources []string, out io.Writer) string {
	cursor := 0
	nodes := make([]string, 0, len(sources))
	draws := make([]string, 0, len(sources))
	b := New()
	defer b.Close()
	for _, text := range sources {
		val := truncate(text, b.Width-1)
		if val != "" {
			nodes = append(nodes, val)
			draws = append(draws, val)
		}
	}
	io.WriteString(out, CURSOR_OFF)
	defer io.WriteString(out, CURSOR_ON)

	if len(nodes) <= 0 {
		nodes = []string{""}
		draws = []string{""}
	}

	offset := 0
	for {
		draws[cursor] = BOLD_ON + truncate(nodes[cursor], b.Width-1) + BOLD_OFF
		status, _, h := b.PrintNoLastLineFeed(nil, draws, offset, out)
		if !status {
			return ""
		}
		draws[cursor] = truncate(nodes[cursor], b.Width-1)
		last := cursor
		for last == cursor {
			switch b.GetCmd() {
			case LEFT:
				cursor = (cursor + len(nodes) - h) % len(nodes)
			case RIGHT:
				cursor = (cursor + h) % len(nodes)
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
			case ENTER:
				return nodes[cursor]
			case LEAVE:
				return ""
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
			fmt.Fprintf(out, "\x1B[%dA\r", h-1)
		} else {
			fmt.Fprintf(out, "\x1B[%dA\r", b.Height-1)
		}
	}
}
