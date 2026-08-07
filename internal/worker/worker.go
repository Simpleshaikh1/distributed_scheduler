package worker

import (
	"context"
	"github.com/Simpleshaikh1/distributed_scheduler/internal/job"
)

type Executor interface {
	Execute(
		context.Context,
		job.Execution,
	) error
}

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

// Now I make the pool own the workers and the queue
type Pool struct {
	workers []*Worker

	jobs chan job.Execution
}

func (p *Pool) Start(ctx context.Context) {
	for _, worker := range p.workers {
		go worker.Run(ctx)
	}
}
