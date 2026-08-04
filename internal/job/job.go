package job

import (
	"errors"
	"time"
)

var ErrInvalidInterval = errors.New(
	"schedule interval must be greater than zero",
)

type Job struct {
	ID       string
	Name     string
	Schedule Schedule
	Action   Action
}

type Schedule struct {
	Interval time.Duration
}

type Action struct {
	Type   ActionType
	Method string
	URL    string
}

type ActionType string

const (
	ActionHttp ActionType = "http"
)

func (s Schedule) Validate() error {
	if s.Interval <= 0 {
		return ErrInvalidInterval
	}

	return nil
}
