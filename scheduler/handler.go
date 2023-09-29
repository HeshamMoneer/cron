package scheduler

import "time"

// handler returns a Job guarded from errors. It decorates the job by the looper before executing it.
//
// Parameters:
//	expectedDuration (time.Duration): The time expected to be taken for the job execution.
//	period (time.Duration): The period of the job recurrence (i.e., the job repeats once every "period" amount of time).
//	job (instruments.Job): The job code that should be executed periodically.
//	indentifier (int): An identifier used by the scheduler to label the jobs.
//
// Returns:
//
// Job is the job with all decorations as should be stored in the cron job pool.
func (c *cron) handler(expectedDuration time.Duration, period time.Duration, job Job, identifier int) Job {
	handled := func() {
		defer func() {
			if err := recover(); err != nil {
				c.log.Error(err)
				c.log.Warn("Exited job with id", identifier)
			}
		}()

		c.looper(expectedDuration, period, job, identifier)()
	}

	return handled
}
