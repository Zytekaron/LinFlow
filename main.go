package main

// #cgo CXXFLAGS: -I/usr/lib/
// #cgo LDFLAGS: -L/usr/lib/ -lstdc++ -lX11
// #include "kbd.hpp"
import "C"
import (
	"github.com/micmonay/keybd_event"
	"golang.design/x/clipboard"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	KeyLShift    = 0xffe1 + iota // Left Shift
	KeyRShift                    // Right shift
	KeyLCtrl                     // Left control
	KeyRCtrl                     // Right control
	KeyCapsLock                  // Caps lock
	KeyShiftLock                 // Shift lock
	KeyLMeta                     // Left meta
	KeyRMeta                     // Right meta
	KeyLAlt                      // Left alt
	KeyRAlt                      // Right alt
	KeyLSuper                    // Left super
	KeyRSuper                    // Right super
	KeyLHyper                    // Left hyper
	KeyRHyper                    // Right hyper
)

type Key struct {
	Code uint32
	Name string
}

var macros = map[string][]byte{}
var buffer []Key
var modeSet bool
var dataMux sync.Mutex // for buffer, modeSet

var kb keybd_event.KeyBonding

func init() {
	var err error
	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(2 * time.Second)
}

func main() {
	C.start_hook()
}

func runMacro() {
	dataMux.Lock()
	defer dataMux.Unlock()

	if len(buffer) == 0 {
		return
	}

	var buf strings.Builder
	buf.WriteString(buffer[0].Name)
	for i := 1; i < len(buffer); i++ {
		buf.WriteRune('+')
		buf.WriteString(buffer[i].Name)
	}

	macro := buf.String()

	if modeSet {
		println("creating macro:", macro)
		macros[macro] = clipboard.Read(clipboard.FmtText)
		return
	}

	// macro does not exist or is empty
	if len(macros[macro]) == 0 {
		return
	}

	println("running macro", macro, "with text:", string(macros[macro]))

	stored := clipboard.Read(clipboard.FmtText)

	clipboard.Write(clipboard.FmtText, macros[macro])

	simulatePaste()

	clipboard.Write(clipboard.FmtText, stored)

}

//export native_pass_key
func native_pass_key(keySym C.ulong, value *C.char) {
	go func() {
		key := Key{
			Code: uint32(keySym),
			Name: cleanKeyName(C.GoString(value)),
		}
		println("native_pass_key call:", key.Code, key.Name)

		dataMux.Lock()
		defer dataMux.Unlock()
		if (keySym == KeyLCtrl || keySym == KeyRCtrl) && len(buffer) == 0 {
			modeSet = true
			return
		}

		buffer = append(buffer, key)
	}()
}

//export native_handle_buffer
func native_handle_buffer() {
	go func() {
		println("native_handle_buffer call")

		runMacro()

		dataMux.Lock()
		defer dataMux.Unlock()
		buffer = nil
		modeSet = false
	}()
}

func cleanKeyName(name string) string {
	name = strings.Replace(name, "KP_macros[macro]", "", 1)
	return name
}

func simulatePaste() {
	kb.SetKeys(keybd_event.VK_V)
	kb.HasCTRL(true)

	err := kb.Press()
	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(time.Microsecond)

	err = kb.Release()
	if err != nil {
		log.Println(err)
		return
	}
}
