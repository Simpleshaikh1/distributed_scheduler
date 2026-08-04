package job

import "time"

type Execution struct {
	ID          string
	JobID       string
	ScheduledAt time.Time
}

func NewExecutionID(jobID string, scheduledAt time.Time) string {
	return jobID + ":" + scheduledAt.UTC().Format(time.RFC3339Nano)
}

func NewExecution(jobID string, scheduledAt time.Time) Execution {
	return Execution{
		ID:          NewExecutionID(jobID, scheduledAt),
		JobID:       jobID,
		ScheduledAt: scheduledAt,
	}
}
