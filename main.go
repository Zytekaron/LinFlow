package main

// #cgo CXXFLAGS: -I/usr/lib/
// #cgo LDFLAGS: -L/usr/lib/ -lstdc++ -lX11
// #include "kbd.hpp"
import "C"

const (
	KeyLShift    = 0xffe1 + iota // Left Shift
	KeyRShift    = 0xffe2        // Right shift
	KeyLCtrl     = 0xffe3        // Left control
	KeyRCtrl     = 0xffe4        // Right control
	KeyCapsLock  = 0xffe5        // Caps lock
	KeyShiftLock = 0xffe6        // Shift lock
	KeyLMeta     = 0xffe7        // Left meta
	KeyRMeta     = 0xffe8        // Right meta
	KeyLAlt      = 0xffe9        // Left alt
	KeyRAlt      = 0xffea        // Right alt
	KeyLSuper    = 0xffeb        // Left super
	KeyRSuper    = 0xffec        // Right super
	KeyLHyper    = 0xffed        // Left hyper
	KeyRHyper    = 0xffee        // Right hyper
)

type Key struct {
	Code uint32
	Name string
}

var buffer []Key

func main() {
	C.start_hook()
}

//export go_pass_key
func go_pass_key(keySym C.ulong, value *C.char) {
	key := Key{
		Code: uint32(keySym),
		Name: C.GoString(value),
	}

	println("[go] c passed Key:", key.Code, key.Name)

	buffer = append(buffer, key)
}

//export go_handle_buffer
func go_handle_buffer() {
	println("[go] c called handle buffer")

	for _, key := range buffer {
		println(key.Code, "->", key.Name)
	}

	buffer = nil
}
