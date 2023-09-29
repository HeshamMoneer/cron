package scheduler

import "time"

// looper returns a Job that runs periodically. It decorates the job by the timer before executing it.
//
// Parameters:
//	expectedDuration (time.Duration): The time expected to be taken for the job execution.
//	period (time.Duration): The period of the job recurrence (i.e., the job repeats once every "period" amount of time).
//	job (instruments.Job): The job code that should be executed periodically.
//	indentifier (int): An identifier used by the scheduler to label the jobs.
//
// Returns:
//
// Job is the job with the loop and the time decorations.
func (c *cron) looper(expectedDuration time.Duration, period time.Duration, job Job, identifier int) Job {
	loop := func() {
		for c.running[identifier] {
			c.log.Info("Started Job with id", identifier)
			c.timer(expectedDuration, period, job, identifier)()
		}

		c.wg.Done()
	}

	return loop
}
