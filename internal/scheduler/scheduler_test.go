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

func TestCalculateDue_CatchesUp(t *testing.T) {
	base := time.Unix(1_000_000, 0).UTC()

	j := job.Job{
		ID: "job-1",
		Schedule: job.Schedule{
			Interval: 5 * time.Minute,
		},
	}

	result := CalculateDue(
		j,
		base.Add(5*time.Minute),
		base.Add(17*time.Minute),
	)

	if len(result.Executions) != 3 {
		t.Fatalf(
			"expected 3 executions, got %d",
			len(result.Executions),
		)
	}

	expected := []time.Time{
		base.Add(5 * time.Minute),
		base.Add(10 * time.Minute),
		base.Add(15 * time.Minute),
	}

	for i, execution := range result.Executions {
		if !execution.ScheduledAt.Equal(expected[i]) {
			t.Errorf(
				"execution %d: got %v, want %v",
				i,
				execution.ScheduledAt,
				expected[i],
			)
		}
	}

	if !result.NextRunAt.Equal(
		base.Add(20 * time.Minute),
	) {
		t.Fatalf(
			"unexpected NextRunAt: %v",
			result.NextRunAt,
		)
	}
}
