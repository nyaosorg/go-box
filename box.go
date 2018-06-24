package box

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

var ansiCutter = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

// Print outputs the list of `nodes` using screen width all
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
	iEnd := len(lines)
	if b.Height > 0 {
		if iEnd >= offset+b.Height {
			iEnd = offset + b.Height
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
				fmt.Fprint(out, escEraseLine)
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
		if i >= iEnd {
			break
		}
		fmt.Fprintln(out)
	}
	return true, nodePerLine, nlines
}

const (
	escCursorOff = "\x1B[?25l"
	escCursorOn  = "\x1B[?25h"
	escBoldOn    = "\x1B[0;47;30m"
	escBoldOff   = "\x1B[0m"
	escUpNChar   = "\x1B[%dA"
	escEraseLine = "\x1B[0K"

	keyLeft  = 0x25
	keyRight = 0x27
	keyUp    = 0x26
	keyDown  = 0x28
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
	io.WriteString(out, escCursorOff)
	defer io.WriteString(out, escCursorOn)

	if len(nodes) <= 0 {
		nodes = []string{""}
		draws = []string{""}
	}

	offset := 0
	for {
		draws[cursor] = escBoldOn + truncate(nodes[cursor], b.Width-2) + escBoldOff
		status, _, h := b.PrintNoLastLineFeed(nil, draws, offset, out)
		if !status {
			return ""
		}
		draws[cursor] = truncate(nodes[cursor], b.Width-2)
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
			fmt.Fprintf(out, escUpNChar+"\r", h-1)
		} else {
			fmt.Fprintf(out, escUpNChar+"\r", b.Height-1)
		}
	}
}
