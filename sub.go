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

type BoxT struct {
	width  int
	height int
	cache  [][]byte
	tty    *tty.TTY
}

func New() *BoxT {
	tty1, err := tty.Open()
	if err != nil {
		panic(err)
	}
	w, h, err := tty1.Size()
	return &BoxT{
		width:  w,
		height: h,
		tty:    tty1,
	}
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
