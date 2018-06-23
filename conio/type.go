package conio

import (
	"syscall"
)

// Handle is the alias of syscall.Handle
type Handle = syscall.Handle

var kernel32 = syscall.NewLazyDLL("kernel32")
