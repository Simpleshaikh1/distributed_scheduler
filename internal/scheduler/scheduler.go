package scheduler

import (
	"context"
	"github.com/Simpleshaikh1/distributed_scheduler/internal/job"
	"sync"
	"time"
)

type Scheduler struct {
	mu sync.Mutex

	jobs map[string]*ScheduledJob

	wakeup chan struct{}
}

type ScheduledJob struct {
	Job       job.Job
	NextRunAt time.Time
}

func DueExecutions(j job.Job, nextRunAt time.Time, now time.Time) []job.Execution {
	var executions []job.Execution

	for !nextRunAt.After(now) {
		executions = append(executions, job.NewExecution(j.ID, nextRunAt))

		nextRunAt = nextRunAt.Add(j.Schedule.Interval)
	}

	return executions
}

type DueResult struct {
	Executions []job.Execution
	NextRunAt  time.Time
}

func CalculateDue(j job.Job, nextRunAt time.Time, now time.Time) DueResult {
	result := DueResult{
		NextRunAt: nextRunAt,
	}

	for !nextRunAt.After(now) {
		result.Executions = append(result.Executions, job.NewExecution(
			j.ID,
			result.NextRunAt,
		))

		result.NextRunAt = result.NextRunAt.Add(j.Schedule.Interval)
	}

	return result
}
