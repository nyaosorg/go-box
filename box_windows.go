// +build windows

package box

import (
	getch "github.com/zetamatta/go-getch"
)

type box_t struct {
	Width  int
	Height int
}

func New() *box_t {
	w, h := GetScreenBufferInfo().ViewSize()
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
	case K_LEFT:
		return LEFT
	case K_RIGHT:
		return RIGHT
	case K_DOWN:
		return DOWN
	case K_UP:
		return UP
	}
	return NONE
}

func (b *box_t) Close() {
}
