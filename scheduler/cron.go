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
}

func Cron() cron {
	return cron{
		wg:   sync.WaitGroup{},
		jobs: map[int]inst.Job{},
	}
}

func (c *cron) AddFunction(expectedDuration time.Duration, period time.Duration, job inst.Job, identifier int) {
	c.wg.Add(1)
	f := func() {
		for {
			job()
			time.Sleep(period)
		}
	}
	c.jobs[identifier] = f
}

func (c *cron) RunAll() {
	for _, job := range c.jobs {
		go job()
	}

	fmt.Scanln()
}
