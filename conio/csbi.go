// +build windows

package conio

import (
	"unsafe"
)

type coordT struct {
	x int16
	y int16
}

func (c coordT) X() int         { return int(c.x) }
func (c coordT) Y() int         { return int(c.y) }
func (c coordT) XY() (int, int) { return int(c.x), int(c.y) }

type smallRectT struct {
	left   int16
	top    int16
	right  int16
	bottom int16
}

func (s smallRectT) Left() int   { return int(s.left) }
func (s smallRectT) Top() int    { return int(s.top) }
func (s smallRectT) Right() int  { return int(s.right) }
func (s smallRectT) Bottom() int { return int(s.bottom) }

// ConsoleScreenBufferInfoT is the type for structure contains terminal's information.
type ConsoleScreenBufferInfoT struct {
	Size              coordT
	CursorPosition    coordT
	Attributes        uint16
	Window            smallRectT
	MaximumWindowSize coordT
}

var getConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")

// GetConsoleScreenBufferInfo returns the latest ConsoleScreenBufferInfoT
// cursor position, window region.
func GetConsoleScreenBufferInfo() *ConsoleScreenBufferInfoT {
	var csbi ConsoleScreenBufferInfoT
	getConsoleScreenBufferInfo.Call(
		uintptr(ConOut()),
		uintptr(unsafe.Pointer(&csbi)))
	return &csbi
}

// ViewSize returns window size from ConsoleScreenBufferInfo structure.
func (csbi *ConsoleScreenBufferInfoT) ViewSize() (int, int) {
	return csbi.Window.Right() - csbi.Window.Left() + 1,
		csbi.Window.Bottom() - csbi.Window.Top() + 1
}
