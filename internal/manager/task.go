package manager

import "context"

type Task interface {
	Start(context.Context)
	GetDescription() string
}

// ErrorTask is an optional extension for tasks that can report a terminal
// failure after Start returns. Task deliberately remains unchanged so that
// existing task implementations and plugins keep their current contract.
type ErrorTask interface {
	Task
	Err() error
}

func taskError(task Task) error {
	if errorTask, ok := task.(ErrorTask); ok {
		return errorTask.Err()
	}
	return nil
}
