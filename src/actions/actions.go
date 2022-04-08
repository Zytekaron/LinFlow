package actions

import (
	"errors"
	"strings"
)

var modKeys = map[string]bool{}

func init() {
	keys := []string{
		"ctrl", "control", "lctrl", "lcontrol", "rctrl", "rcontrol",
		"cmd", "command", "lcmd", "lcommand", "rcmd", "rcommand",
		"alt", "lalt", "ralt",
	}
	for _, key := range keys {
		modKeys[key] = true
	}
}

type ActionType string

const (
	ActionText    ActionType = "TEXT"
	ActionWait    ActionType = "WAIT"
	ActionMove    ActionType = "MOVE"
	ActionClick   ActionType = "CLICK"
	ActionKeybind ActionType = "KEYBIND"
)

type Action interface {
	Execute() error
}

func Parse(input string) (Action, error) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) == 0 {
		return nil, errors.New("malformed command: missing action")
	}

	var parseInput string
	if len(parts) == 2 {
		parseInput = parts[1]
	}

	switch ActionType(parts[0]) {
	case ActionText:
		return parseText(parseInput)
	case ActionWait:
		return parseWait(parseInput)
	case ActionMove:
		return parseMove(parseInput)
	case ActionClick:
		return parseClick(parseInput)
	case ActionKeybind:
		return parseKeybind(parseInput)
	default:
		return nil, errors.New("invalid action type")
	}
}
