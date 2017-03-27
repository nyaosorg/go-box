package box

import (
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

type box_t struct {
	Width  int
	Height int
}

func (b *box_t) Print(ctx context.Context,
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

	for i := offset; i < i_end; i++ {
		fmt.Fprintln(out, string(lines[i]))
		if ctx != nil {
			select {
			case <-ctx.Done():
				return false, nodePerLine, nlines
			default:
			}
		}
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
	for _, text := range sources {
		val := truncate(text, b.Width-1)
		if val != "" {
			nodes = append(nodes, val)
			draws = append(draws, val)
		}
	}
	io.WriteString(out, CURSOR_OFF)
	defer io.WriteString(out, CURSOR_ON)

	offset := 0
	for {
		draws[cursor] = BOLD_ON + truncate(nodes[cursor], b.Width-1) + BOLD_OFF
		status, _, h := b.Print(nil, draws, offset, out)
		if !status {
			return ""
		}
		draws[cursor] = truncate(nodes[cursor], b.Width-1)
		last := cursor
		for last == cursor {
			switch get() {
			case LEFT:
				if cursor-h >= 0 {
					cursor -= h
				}
			case RIGHT:
				if cursor+h < len(nodes) {
					cursor += h
				}
			case DOWN:
				if cursor+1 < len(nodes) {
					cursor++
				}
			case UP:
				if cursor > 0 {
					cursor--
				}
			case ENTER:
				return nodes[cursor]
			case LEAVE:
				return ""
			}

			// x := cursor / h
			y := cursor % h
			if y < offset {
				offset--
			} else if y >= offset+b.Height {
				offset++
			}
		}
		if h < b.Height {
			fmt.Fprintf(out, "\x1B[%dA", h)
		} else {
			fmt.Fprintf(out, "\x1B[%dA", b.Height)
		}
	}
}
