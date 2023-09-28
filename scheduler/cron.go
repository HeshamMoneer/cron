package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type Cron struct {
	wg   sync.WaitGroup
	jobs []Job
}

func CronInit() Cron {
	return Cron{
		wg:   sync.WaitGroup{},
		jobs: make([]Job, 0),
	}
}

func (c *Cron) AddFunction(period time.Duration, job Job) {
	c.wg.Add(1)
	f := func() {
		for {
			job()
			time.Sleep(period)
		}
	}
	c.jobs = append(c.jobs, f)
}

func (c *Cron) RunAll() {
	for _, job := range c.jobs {
		go job()
	}

	fmt.Scanln()
}
