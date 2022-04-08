package actions

import (
	"errors"
	"github.com/go-vgo/robotgo"
	"strings"
)

type KeybindAction struct {
	Key  string
	Mods []string
}

func (k *KeybindAction) Execute() error {
	if len(k.Mods) > 0 {
		robotgo.KeyTap(k.Key, k.Mods)
	} else {
		robotgo.KeyTap(k.Key)
	}
	return nil
}

func parseKeybind(input string) (*KeybindAction, error) {
	if len(input) == 0 {
		return nil, errors.New("not enough arguments for keybind action")
	}

	input = strings.ToLower(input)
	keys := strings.Split(input, "+")

	var key string
	var mods []string
	for _, k := range keys {
		switch {
		case modKeys[k]:
			mods = append(mods, k)
		case key == "":
			key = k
		default:
			return nil, errors.New("only one non-mod key is supported for keybind action")
		}
	}

	return &KeybindAction{
		Key:  key,
		Mods: mods,
	}, nil
}
