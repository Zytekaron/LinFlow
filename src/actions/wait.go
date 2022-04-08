package actions

import (
	"errors"
	"github.com/cstockton/go-conv"
	"time"
)

type WaitAction struct {
	Duration time.Duration
}

func (w *WaitAction) Execute() error {
	time.Sleep(w.Duration)
	return nil
}

func parseWait(input string) (*WaitAction, error) {
	duration, err := conv.Duration(input)
	if err != nil {
		return nil, errors.New("could not parse duration from input")
	}

	return &WaitAction{
		Duration: duration,
	}, nil
}
