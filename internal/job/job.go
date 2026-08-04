package job

import "time"

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
