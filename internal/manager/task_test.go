package manager

import (
	"context"
	"errors"
	"testing"
)

type errorReportingTestTask struct {
	err error
}

func (t *errorReportingTestTask) Start(context.Context) {}

func (t *errorReportingTestTask) GetDescription() string {
	return "error-reporting test task"
}

func (t *errorReportingTestTask) Err() error {
	return t.err
}

type nonReportingTestTask struct{}

func (nonReportingTestTask) Start(context.Context) {}

func (nonReportingTestTask) GetDescription() string {
	return "non-reporting test task"
}

func TestTaskError(t *testing.T) {
	expected := errors.New("sprite generation failed")

	tests := []struct {
		name string
		task Task
		want error
	}{
		{
			name: "reports failure",
			task: &errorReportingTestTask{err: expected},
			want: expected,
		},
		{
			name: "reports success",
			task: &errorReportingTestTask{},
		},
		{
			name: "preserves non-reporting task compatibility",
			task: nonReportingTestTask{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := taskError(tt.task); !errors.Is(got, tt.want) {
				t.Fatalf("taskError() = %v, want %v", got, tt.want)
			}
		})
	}
}
