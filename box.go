package box

import (
	"bufio"
	"fmt"
	"io"
	"sort"

	"github.com/nyaosorg/go-ttyadapter"
	"github.com/nyaosorg/go-ttyadapter/tty8"

	"github.com/nyaosorg/go-box/v3/grid"

	"github.com/nyaosorg/go-box/v3/internal/ansi"
	"github.com/nyaosorg/go-box/v3/internal/keys"
)

type Box struct {
	grid.Grid
	ttyadapter.Tty
}

// New creates and initializes a Box using the default terminal backend (tty8).
// It is equivalent to creating a Box with &tty8.Tty{} and calling Open().
//
// Example:
//
//	b, err := box.New()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer b.Close()
//
// Use this function when you do not need to customize the terminal input.
func New() (*Box, error) {
	b := &Box{
		Tty: &tty8.Tty{},
	}
	return b, b.Open()
}

// Open initializes the Box with the terminal backend assigned to b.Tty.
// This method is used when you want to provide a custom TTY implementation
// instead of the default one created by New().
//
// Example:
//
//	import "github.com/nyaosorg/go-ttyadapter/auto"
//
//	b := &box.Box{
//		Tty: &auto.Pilot{Text: []string{"l", "l", "\r"}},
//	}
//	if err := b.Open(); err != nil {
//		log.Fatal(err)
//	}
//	defer b.Close()
func (b *Box) Open() error {
	err := b.Tty.Open(nil)
	if err != nil {
		return err
	}
	b.Width, b.Height, err = b.Tty.Size()
	return err
}

// SelectIndex returns the indexes that user selected.
func (b *Box) SelectIndex(sources []string, multi bool, out io.Writer) ([]int, error) {
	cursor := 0
	selected := make(map[int]struct{})

	nodes := make([]nodeT, 0, len(sources))
	draws := make([]string, 0, len(sources))
	for i, text := range sources {
		val := truncate(text, b.Width-1)
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
			draws[index] = ansi.BoldOn + truncate(nodes[index].Text, b.Width-2) + ansi.BoldOff
		}
		draws[cursor] = ansi.BoldOn2 + truncate(nodes[cursor].Text, b.Width-2) + ansi.BoldOff
		_, h, err := b.Print(draws, offset, out)
		if err != nil {
			return []int{}, err
		}
		for index := range selected {
			draws[index] = truncate(nodes[index].Text, b.Width-2)
		}
		draws[cursor] = truncate(nodes[cursor].Text, b.Width-2)
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
			key, err := b.Tty.GetKey()
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
			} else if y >= offset+b.Height {
				offset = y - b.Height + 1
				// offset++
			}
		}
		if h < b.Height {
			if h > 1 {
				fmt.Fprintf(out, ansi.UpN, h-1)
			}
		} else {
			if b.Height > 1 {
				fmt.Fprintf(out, ansi.UpN, b.Height-1)
			}
		}
		fmt.Fprint(out, "\r")
	}
}

// SelectString returns the strings that user selected.
func (b *Box) SelectString(sources []string, multi bool, out io.Writer) ([]string, error) {
	list, err := b.SelectIndex(sources, multi, out)
	if err != nil {
		return nil, err
	}
	values := make([]string, 0, len(list))
	for _, index := range list {
		values = append(values, sources[index])
	}
	return values, nil
}
