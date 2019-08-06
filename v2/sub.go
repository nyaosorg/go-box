package box

import (
	"strings"

	"github.com/mattn/go-tty"
)

const (
	K_LEFT       = "\x1B[D"
	K_UP         = "\x1B[A"
	K_RIGHT      = "\x1B[C"
	K_DOWN       = "\x1B[B"
	K_CTRL_F     = "\x06"
	K_CTRL_B     = "\x02"
	K_CTRL_N     = "\x0E"
	K_CTRL_P     = "\x10"
	K_CTRL_G     = "\x07"
	K_CTRL_DOWN  = "\x1B[1;5B"
	K_CTRL_LEFT  = "\x1B[1;5D"
	K_CTRL_RIGHT = "\x1B[1;5C"
	K_CTRL_UP    = "\x1B[1;5A"
)

type box_t struct {
	Width  int
	Height int
	Cache  [][]byte
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

func (b *box_t) getkey() (string, error) {
	var keys strings.Builder
	clean, err := b.Tty.Raw()
	if err != nil {
		return "", err
	}
	defer clean()
	for {
		key, err := b.Tty.ReadRune()
		if err != nil {
			return "", err
		}
		if key == 0 {
			continue
		}
		keys.WriteRune(key)
		if !b.Tty.Buffered() {
			return keys.String(), nil
		}
	}
}

func (b *box_t) GetCmd() int {
	key, err := b.getkey()
	if err != nil {
		return NONE
	}
	switch key {
	case "h", K_CTRL_B, K_LEFT:
		return LEFT
	case "l", K_CTRL_F, K_RIGHT:
		return RIGHT
	case "j", K_CTRL_N, K_DOWN:
		return DOWN
	case "k", K_CTRL_P, K_UP:
		return UP
	case " ", "J", K_CTRL_DOWN:
		return SELECT_DOWN
	case "\b", "K", K_CTRL_UP:
		return SELECT_UP
	case "H", K_CTRL_LEFT:
		return SELECT_LEFT
	case "L", K_CTRL_RIGHT:
		return SELECT_RIGHT
	case "\r", "\n":
		return ENTER
	case "\x1B", K_CTRL_G:
		return LEAVE
	}
	return NONE
}

func (b *box_t) Close() {
	b.Tty.Close()
}
