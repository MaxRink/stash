package manager

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
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

type blockingErrorTask struct {
	started chan<- struct{}
	release <-chan struct{}
	count   *int32
	err     error
}

func (t *blockingErrorTask) Start(context.Context) {
	t.started <- struct{}{}
	atomic.AddInt32(t.count, 1)
	<-t.release
}

func (t *blockingErrorTask) GetDescription() string { return "blocking task" }
func (t *blockingErrorTask) Err() error             { return t.err }

func TestExecuteGenerateTasksWaitsAndAggregates(t *testing.T) {
	started := make(chan struct{}, 3)
	var startedCount int32
	release := make(chan struct{})
	wantErr := errors.New("generation failed")
	queue := make(chan Task, 3)
	for _, err := range []error{wantErr, nil, wantErr} {
		queue <- &blockingErrorTask{started: started, release: release, count: &startedCount, err: err}
	}
	close(queue)

	result := make(chan error, 1)
	go func() {
		result <- executeGenerateTasks(context.Background(), 2, queue, func(task Task) {
			task.Start(context.Background())
		}, nil)
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for bounded tasks to start")
		}
	}
	select {
	case <-started:
		t.Fatal("started more tasks than the concurrency limit")
	case <-time.After(50 * time.Millisecond):
	}

	close(release)
	select {
	case err := <-result:
		if !errors.Is(err, wantErr) {
			t.Fatalf("executeGenerateTasks() = %v, want %v", err, wantErr)
		}
	case <-time.After(time.Second):
		t.Fatal("executeGenerateTasks did not wait for all tasks")
	}
	if got := atomic.LoadInt32(&startedCount); got != 3 {
		t.Fatalf("started %d tasks, want 3", got)
	}
}
