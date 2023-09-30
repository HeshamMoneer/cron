package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type cron struct {
	wg      sync.WaitGroup
	jobs    map[int]Job
	running map[int]bool
	log     Logger
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
		wg:      sync.WaitGroup{},
		jobs:    map[int]Job{},
		running: map[int]bool{},
		log:     NewLogger(file...),
	}
}

// AddJob registers a Job to the cron scheduler job pool to be run periodically when needed.
//
// Parameters:
//
//	expectedDuration (time.Duration): The time expected to be taken for the job execution. This acts as the job timeout value.
//	period (time.Duration): The period of the job recurrence (i.e., the job repeats once every "period" amount of time).
//	job (instruments.Job): The job code that should be executed periodically.
//	indentifier (int): An identifier used by the scheduler to label the jobs.
func (c *cron) AddJob(expectedDuration time.Duration, period time.Duration, job Job, identifier int) {
	c.log.Info("Attempting to register job with id", identifier, "......................")
	if c.running[identifier] {
		c.log.Error("Job with id", identifier, "is already registered and running. Please stop it first then attempt to replace it.")
		return
	}

	c.wg.Add(1)
	c.jobs[identifier] = c.handler(expectedDuration, period, job, identifier)
	c.log.Info("Registered job with id", identifier, "successfully!!!")
}

// RunAll runs all the Jobs registered to the cron job pool.
func (c *cron) RunAll() {
	for identifier := range c.jobs {
		if !c.running[identifier] {
			c.RunJob(identifier)
		}
	}
}

// StopAll stops all the Jobs that are running.
func (c *cron) StopAll() {
	for identifier := range c.jobs {
		if c.running[identifier] {
			c.StopJob(identifier)
		}
	}
}

// RunJob runs a specfic job from cron job pool.
//
// Parameters:
//
//	identifier (int): the cron id of the job that should start running periodically.
func (c *cron) RunJob(identifier int) {
	c.log.Info("Attempting to schedule job with id", identifier, "......................")
	if c.running[identifier] {
		c.log.Warn("Job with id", identifier, "is already scheduled. Nothing will change.")
	} else {
		job, ok := c.jobs[identifier]
		if !ok {
			c.log.Warn("Job with id", identifier, "was not registered. No job will be scheduled.")
			return
		}
		c.running[identifier] = true
		go job()
		c.log.Info("Started scheduling job with id", identifier, "successfully!!!")
	}
}

// StopJob stops a specfic job from cron job pool.
//
// Parameters:
//
//	identifier (int): the cron id of the job that should stop running periodically.
func (c *cron) StopJob(identifier int) {
	c.log.Info("Attempting to stop job with id", identifier, "......................")
	if !c.running[identifier] {
		c.log.Warn("Job with id", identifier, "is either stopped or not existing. Nothing will change.")
	} else {
		c.running[identifier] = false
		c.log.Info("Stopped scheduling job with id", identifier, "successfully!!!")
	}
}

// WaitJobs waits until the user click "Enter" to allow for the Jobs to run periodically.
func (c *cron) WaitJobs() {
	fmt.Scanln()
	c.log.Info("Exiting cron.......")
}
