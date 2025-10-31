package ansi

import (
	"regexp"
)

const (
	CursorOff = "\x1B[?25l"
	CursorOn  = "\x1B[?25h"
	BoldOn    = "\x1B[0;47;30m"
	BoldOn2   = "\x1B[0;1;7m"
	BoldOff   = "\x1B[0m"
	UpN       = "\x1B[%dA"
	EraseLine = "\x1B[0K"
)

var RxSequence = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")
