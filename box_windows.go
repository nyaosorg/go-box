// +build windows

package box

import (
	"github.com/zetamatta/go-box/conio"
	getch "github.com/zetamatta/go-getch"
)

type box_t struct {
	Width  int
	Height int
	Cache  [][]byte
}

// New is the constructor for box_t
func New() *box_t {
	w, h := conio.GetConsoleScreenBufferInfo().ViewSize()
	return &box_t{
		Width:  w,
		Height: h,
	}
}

func (b *box_t) GetCmd() int {
	k := getch.All().Key
	if k == nil {
		return NONE
	}
	switch k.Rune {
	case 'h', ('b' & 0x1F):
		return LEFT
	case 'l', ('f' & 0x1F):
		return RIGHT
	case 'j', ('n' & 0x1F), ' ':
		return DOWN
	case 'k', ('p' & 0x1F), '\b':
		return UP
	case '\r', '\n':
		return ENTER
	case '\x1B', ('g' & 0x1F):
		return LEAVE
	}

	switch k.Scan {
	case keyLeft:
		return LEFT
	case keyRight:
		return RIGHT
	case keyDown:
		return DOWN
	case keyUp:
		return UP
	}
	return NONE
}

func (b *box_t) Close() {
}
