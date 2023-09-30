package scheduler

import (
	"cron/scheduler/timeout"
	"time"
)

// timer returns a Job that is timed (i.e., its start and end times are measured and logged, and the goroutine is sent to sleep for the appropriate amount of time). It executes the actual job.
//
// Parameters:
//
//	expectedDuration (time.Duration): The time expected to be taken for the job execution. This acts as the job timeout value.
//	period (time.Duration): The period of the job recurrence (i.e., the job repeats once every "period" amount of time).
//	job (instruments.Job): The job code that should be executed periodically.
//	indentifier (int): An identifier used by the scheduler to label the jobs.
//
// Returns:
//
// Job is the job with the time decoration.
func (c *cron) timer(expectedDuration time.Duration, period time.Duration, job Job, identifier int) Job {
	timed := func() {
		startTime := time.Now()

		wrapped := func() error {
			job()
			return nil
		}

		err := timeout.Run(wrapped, expectedDuration)
		if err != nil {
			panic(err)
		}

		endTime := time.Now()

		actualDuration := endTime.Sub(startTime)
		timeToSleep := period - actualDuration
		if timeToSleep < 0 {
			timeToSleep = 0
		}
		c.log.Info("Finished Job with id", identifier, "Timeout Duration:", expectedDuration, "Actual Duration:", actualDuration)
		time.Sleep(timeToSleep)
	}

	return timed
}
