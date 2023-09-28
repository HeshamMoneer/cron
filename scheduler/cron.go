package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type Cron struct {
	wg   sync.WaitGroup
	jobs map[int]Job
}

func CronInit() Cron {
	return Cron{
		wg:   sync.WaitGroup{},
		jobs: map[int]Job{},
	}
}

func (c *Cron) AddFunction(expectedDuration time.Duration, period time.Duration, job Job, identifier int) {
	c.wg.Add(1)
	f := func() {
		for {
			job()
			time.Sleep(period)
		}
	}
	c.jobs[identifier] = f
}

func (c *Cron) RunAll() {
	for _, job := range c.jobs {
		go job()
	}

	fmt.Scanln()
}
