package box

import (
	"strings"

	"github.com/mattn/go-tty"
)

const (
	_K_LEFT       = "\x1B[D"
	_K_UP         = "\x1B[A"
	_K_RIGHT      = "\x1B[C"
	_K_DOWN       = "\x1B[B"
	_K_CTRL_F     = "\x06"
	_K_CTRL_B     = "\x02"
	_K_CTRL_N     = "\x0E"
	_K_CTRL_P     = "\x10"
	_K_CTRL_G     = "\x07"
	_K_CTRL_DOWN  = "\x1B[1;5B"
	_K_CTRL_LEFT  = "\x1B[1;5D"
	_K_CTRL_RIGHT = "\x1B[1;5C"
	_K_CTRL_UP    = "\x1B[1;5A"
	_K_SHIFT_TAB  = "\x1B[Z"
)

type BoxT struct {
	width  int
	height int
	cache  [][]byte
	tty    *tty.TTY
}

func NewBox() (*BoxT, error) {
	tty1, err := tty.Open()
	if err != nil {
		return nil, err
	}
	w, h, err := tty1.Size()
	return &BoxT{
		width:  w,
		height: h,
		tty:    tty1,
	}, err
}

func (b *BoxT) getKey() (string, error) {
	var keys strings.Builder
	clean, err := b.tty.Raw()
	if err != nil {
		return "", err
	}
	defer clean()
	for {
		key, err := b.tty.ReadRune()
		if err != nil {
			return "", err
		}
		if key == 0 {
			continue
		}
		keys.WriteRune(key)
		if !b.tty.Buffered() {
			return keys.String(), nil
		}
	}
}

func (b *BoxT) Close() {
	b.tty.Close()
}
