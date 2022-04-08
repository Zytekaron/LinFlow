package actions

import (
	"errors"
	"github.com/go-vgo/robotgo"
	"strconv"
	"strings"
)

type MoveAction struct {
	X int
	Y int
}

func (m *MoveAction) Execute() error {
	robotgo.Move(m.X, m.Y)
	return nil
}

func parseMove(input string) (*MoveAction, error) {
	partsXY := strings.Split(input, " ")
	if len(partsXY) < 2 {
		return nil, errors.New("not enough arguments for move action")
	}

	x, err := strconv.Atoi(partsXY[1])
	if err != nil {
		return nil, errors.New("could not parse x")
	}

	y, err := strconv.Atoi(partsXY[2])
	if err != nil {
		return nil, errors.New("could not parse y")
	}

	return &MoveAction{X: x, Y: y}, nil
}
