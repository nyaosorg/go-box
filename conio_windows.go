// +build windows

package box

import (
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32")

var hConout syscall.Handle

func init() {
	var err error
	hConout, err = syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		panic(err.Error())
	}
}

type coord_t struct {
	X int16
	Y int16
}

type small_rect_t struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type console_screen_buffer_info_t struct {
	Size              coord_t
	CursorPosition    coord_t
	Attributes        uint16
	Window            small_rect_t
	MaximumWindowSize coord_t
}

var getConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")

type console_handle_t syscall.Handle

func newHandle(handle syscall.Handle) console_handle_t {
	return console_handle_t(handle)
}

func (h console_handle_t) GetScreenBufferInfo() *console_screen_buffer_info_t {
	var csbi console_screen_buffer_info_t
	getConsoleScreenBufferInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&csbi)))
	return &csbi
}

func GetScreenBufferInfo() *console_screen_buffer_info_t {
	return console_handle_t(hConout).GetScreenBufferInfo()
}

func (this *console_screen_buffer_info_t) ViewSize() (int, int) {
	return int(this.Window.Right-this.Window.Left) + 1,
		int(this.Window.Bottom-this.Window.Top) + 1
}

func (this *console_screen_buffer_info_t) CursorPos() (int, int) {
	return int(this.CursorPosition.X), int(this.CursorPosition.Y)
}

func GetLocate() (int, int) {
	return GetScreenBufferInfo().CursorPos()
}
