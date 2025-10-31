package box

import (
	"strings"

	"github.com/mattn/go-tty"
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
