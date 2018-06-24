package conio

import (
	"syscall"
)

type Handle = syscall.Handle

var kernel32 = syscall.NewLazyDLL("kernel32")
