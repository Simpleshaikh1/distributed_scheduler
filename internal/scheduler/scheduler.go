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

func New() *Scheduler {
	return &Scheduler{
		jobs:   make(map[string]*ScheduledJob),
		wakeup: make(chan struc{}, 1),
	}
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

func (s *Scheduler) Add(
	j job.Job,
	firstRunAt time.Time,
) error {
	if err := j.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs[j.ID] = &ScheduledJob{
		Job:       j,
		NextRunAt: firstRunAt.UTC(),
	}

	s.signalWakeup()

	return nil
}

func (s *Scheduler) Remove(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.jobs, id)

	s.signalWakeup()
}

func (s *Scheduler) nextRunAt() (time.Time, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var next time.Time
	found := false

	for _, scheduled := range s.jobs {
		if !found || scheduled.NextRunAt.Before(next) {
			next = scheduled.NextRunAt
			found = true
		}
	}

	return next, found
}

func (s *Scheduler) Run(ctx context.Context) error {
	for {
		nextRunAt, ok := s.nextRunAt()

		if !ok {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-s.wakeup:
				continue
			}
		}

		delay := time.Until(nextRunAt)

		if delay <= 0 {
			s.runDueJobs(time.Now())
			continue
		}

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}

			return ctx.Err()

		case <-s.wakeup:
			if !timer.Stop() {
				<-timer.C
			}

			continue

		case <-timer.C:
			s.runDueJobs(time.Now())
		}
	}
}

func (s *Scheduler) runDueJobs(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, scheduled := range s.jobs {
		result := CalculateDue(
			scheduled.Job,
			scheduled.NextRunAt,
			now,
		)

		scheduled.NextRunAt = result.NextRunAt

		for _, execution := range result.Executions {
			// dispatch later
		}
	}
}

func (s *Scheduler) collectDue(
	now time.Time,
) []job.Execution {
	s.mu.Lock()
	defer s.mu.Unlock()

	var executions []job.Execution

	for _, scheduled := range s.jobs {
		result := CalculateDue(
			scheduled.Job,
			scheduled.NextRunAt,
			now,
		)

		scheduled.NextRunAt = result.NextRunAt

		executions = append(
			executions,
			result.Executions...,
		)
	}

	return executions
}
