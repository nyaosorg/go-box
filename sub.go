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

type _Tty interface {
	GetKey() (string, error)
	Close() error
}

type _GoTty struct {
	*tty.TTY
}

func (g _GoTty) GetKey() (string, error) {
	var keys strings.Builder
	clean, err := g.TTY.Raw()
	if err != nil {
		return "", err
	}
	defer clean()
	for {
		key, err := g.TTY.ReadRune()
		if err != nil {
			return "", err
		}
		if key == 0 {
			continue
		}
		keys.WriteRune(key)
		if !g.TTY.Buffered() {
			return keys.String(), nil
		}
	}
}
