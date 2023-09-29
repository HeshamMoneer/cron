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

func NewCron(file ...string) cron {
	return cron{
		wg:   sync.WaitGroup{},
		jobs: map[int]inst.Job{},
		log:  inst.NewLogger(file...),
	}
}

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

func (c *cron) RunAll() {
	for identifier, job := range c.jobs {
		go inst.HandleErrors(job, identifier, c.log)
	}

	fmt.Scanln()
}
