package worker

import (
	"context"
	"github.com/Simpleshaikh1/distributed_scheduler/internal/job"
)

type Worker struct {
	id int

	jobs <-chan job.Execution

	executor Executor
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case execution, ok := <-w.jobs:
			if !ok {
				return
			}
			_ = w.executor.Execute(ctx, execution)
		}
	}
}
