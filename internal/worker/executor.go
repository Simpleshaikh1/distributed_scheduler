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
