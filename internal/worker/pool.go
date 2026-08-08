package worker

import (
	"context"
	"github.com/Simpleshaikh1/distributed_scheduler/internal/job"
	"sync"
)

// Now I make the pool own the workers and the queue
type Pool struct {
	jobs chan job.Execution

	workers int

	executor Executor

	wg sync.WaitGroup
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)

		worker := Worker{
			id:       1,
			jobs:     p.jobs,
			executor: p.executor,
		}

		go worker.Run(ctx, &p.wg)
	}
}
