package actions

import (
	"errors"
	"github.com/go-vgo/robotgo"
)

type ClickAction struct {
	Button string
}

func (c *ClickAction) Execute() error {
	robotgo.Click(c.Button)
	return nil
}

func parseClick(input string) (*ClickAction, error) {
	if robotgo.MouseMap[input] == 0 {
		return nil, errors.New("unknown mouse button")
	}

	return &ClickAction{
		Button: input,
	}, nil
}
