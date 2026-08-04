package scheduler

import (
	"github.com/Simpleshaikh1/distributed_scheduler/internal/job"
	"testing"
	"time"
)

func TestCalculateDue_NoExecution(t *testing.T) {
	base := time.Unix(1_000_000, 0).UTC()

	j := job.Job{
		ID: "job-1",
		Schedule: job.Schedule{
			Interval: time.Minute,
		},
	}

	result := CalculateDue(
		j,
		base.Add(time.Minute),
		base,
	)

	if len(result.Executions) != 0 {
		t.Fatalf(
			"expected 0 executions, got %d",
			len(result.Executions),
		)
	}

	if !result.NextRunAt.Equal(base.Add(time.Minute)) {
		t.Fatalf(
			"unexpected NextRunAt: %v",
			result.NextRunAt,
		)
	}
}

func TestCalculateDue_OneExecution(t *testing.T) {
	base := time.Unix(1_000_000, 0).UTC()

	j := job.Job{
		ID: "job-1",
		Schedule: job.Schedule{
			Interval: time.Minute,
		},
	}

	result := CalculateDue(
		j,
		base,
		base,
	)

	if len(result.Executions) != 1 {
		t.Fatalf(
			"expected 1 execution, got %d",
			len(result.Executions),
		)
	}

	if !result.Executions[0].ScheduledAt.Equal(base) {
		t.Fatalf(
			"unexpected ScheduledAt: %v",
			result.Executions[0].ScheduledAt,
		)
	}

	if !result.NextRunAt.Equal(
		base.Add(time.Minute),
	) {
		t.Fatalf(
			"unexpected NextRunAt: %v",
			result.NextRunAt,
		)
	}
}
