package actions

import (
	"github.com/go-vgo/robotgo"
)

type TextAction struct {
	Text string
}

func (t *TextAction) Execute() error {
	robotgo.TypeStr(t.Text)
	return nil
}

func parseText(input string) (*TextAction, error) {
	return &TextAction{
		Text: input,
	}, nil
}
