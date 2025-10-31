package ansi

import (
	"regexp"
)

var RxSequence = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")
