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
		wg: sync.WaitGroup{},
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
	for i := 0; i < len(c.jobs); i++ {
		go c.jobs[i]()
	}

	c.wg.Wait()
	fmt.Scanln()
}
