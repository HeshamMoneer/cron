package scheduler

import (
	"fmt"
	"sync"
	"time"

	inst "cron/instruments"
)

type cron struct {
	wg   sync.WaitGroup
	jobs map[int]inst.Job
	log  inst.Logger
}

// NewCron creates a cron scheduler.
//
// Parameters:
//
//	fileName [optional] (string): If passed, the scheduler logs go to a file in the logs folder. Otherwise, the logs go by default to the standard output.
//
// Returns:
//
//	cron: The cron scheduler instance.
func NewCron(file ...string) cron {
	return cron{
		wg:   sync.WaitGroup{},
		jobs: map[int]inst.Job{},
		log:  inst.NewLogger(file...),
	}
}

// AddJob registers a Job to the cron scheduler job pool to be run periodically when needed.
//
// Parameters:
//
//	expectedDuration (time.Duration): The time expected to be taken for the job execution.
//	period (time.Duration): The period of the job recurrence (i.e., the job repeats once every "period" amount of time).
//	job (instruments.Job): The job code that should be executed periodically.
//	indentifier (int): An identifier used by the scheduler to label the jobs.
func (c *cron) AddJob(expectedDuration time.Duration, period time.Duration, job inst.Job, identifier int) {
	c.wg.Add(1)
	loop := func() {
		for {
			c.log.Info("Started Job", identifier)
			startTime, endTime := inst.JobTime(job)
			actualDuration := endTime.Sub(startTime)
			timeToSleep := period - actualDuration

			time.Sleep(timeToSleep)
			c.log.Info("Finished Job", identifier, "Expected Duration:", expectedDuration, "Actual Duration:", actualDuration)
		}
	}
	c.jobs[identifier] = loop
}

// RunAll runs all the Jobs registered to the cron job pool.
func (c *cron) RunAll() {
	for identifier, job := range c.jobs {
		go inst.HandleErrors(job, identifier, c.log)
	}

	fmt.Scanln()
}
