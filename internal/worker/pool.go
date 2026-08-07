package worker

import (
	"context"
	"github.com/Simpleshaikh1/distributed_scheduler/internal/job"
)

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
