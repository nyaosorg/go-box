// -build windows
package box

import "github.com/mattn/go-tty"

type box_t struct {
	Width  int
	Height int
	Tty    *tty.TTY
}

func New() *box_t {
	tty1, err := tty.Open()
	if err != nil {
		panic(err)
	}
	w, h, err := tty1.Size()
	return &box_t{
		Width:  w,
		Height: h,
		Tty:    tty1,
	}
}

func (b *box_t) GetCmd() int {
	key, err := b.Tty.ReadRune()
	if err != nil {
		return NONE
	}
	switch key {
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
	return NONE
}

func (b *box_t) Close() {
	b.Tty.Close()
}
