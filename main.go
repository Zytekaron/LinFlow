package main

// #cgo CXXFLAGS: -I/usr/lib/
// #cgo LDFLAGS: -L/usr/lib/ -lstdc++ -lX11
// #include "kbd.hpp"
import "C"
import (
	"fmt"
	"linflow/src/config"
	"linflow/src/macro"
	"log"
	"strings"
	"sync"
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

var modKeys = map[C.ulong]bool{
	KeyLShift:    true,
	KeyRShift:    true,
	KeyLCtrl:     true,
	KeyRCtrl:     true,
	KeyCapsLock:  true,
	KeyShiftLock: true,
	KeyLMeta:     true,
	KeyRMeta:     true,
	KeyLAlt:      true,
	KeyRAlt:      true,
	KeyLSuper:    true,
	KeyRSuper:    true,
	KeyLHyper:    true,
	KeyRHyper:    true,
}

type Key struct {
	Code uint32
	Name string
}

func (k *Key) String() string {
	return k.Name
}

var cfg *config.Config
var configDirs = []string{
	"/etc/linflow/config.yml",
	"config.yml",
}

var macros = map[string]*macro.Macro{}
var buffer []Key

//var modeSet bool
var bufferMux sync.Mutex

func init() {
	var err error
	cfg, err = config.Load(configDirs)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	for _, m := range cfg.Macros {
		fmt.Println("Loading", m.Code)

		mac, err := m.ToMacro()
		if err != nil {
			fmt.Println("parse", err)
			return
		}
		macros[mac.Code] = mac
	}

	C.start_hook(C.int(cfg.LinMod))
}

func execute() {
	bufferMux.Lock()
	defer bufferMux.Unlock()

	if len(buffer) == 0 {
		return
	}

	name := keysToString(buffer)
	mac, ok := macros[name]
	if !ok {
		return
	}

	fmt.Println("executing macro:", name)
	err := mac.Execute()
	if err != nil {
		log.Println(err)
	}
}

//export native_pass_key
//goland:noinspection GoSnakeCaseUsage
func native_pass_key(keySym C.ulong, value *C.char) {
	go func() {
		// not accepting mod keys due to shifted keys like '?' is SHIFT+/
		if modKeys[keySym] {
			return
		}

		key := Key{
			Code: uint32(keySym),
			Name: cleanKeyName(C.GoString(value)),
		}
		fmt.Println("native_pass_key call", key.Code, key.Name)

		bufferMux.Lock()
		buffer = append(buffer, key)
		bufferMux.Unlock()
	}()
}

//export native_handle_buffer
//goland:noinspection GoSnakeCaseUsage
func native_handle_buffer() {
	go func() {
		fmt.Println("native_handle_buffer call")

		execute()

		bufferMux.Lock()
		defer bufferMux.Unlock()
		buffer = nil
		//modeSet = false
	}()
}

func cleanKeyName(name string) string {
	return strings.Replace(name, "KP_macros[macro]", "", 1)
}

func keysToString(args []Key) string {
	if len(args) == 0 {
		return ""
	}

	var buf strings.Builder
	buf.WriteString(args[0].String())
	for i := 1; i < len(args); i++ {
		buf.WriteString("+")
		buf.WriteString(args[i].String())
	}

	return buf.String()
}
