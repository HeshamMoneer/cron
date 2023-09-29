package scheduler

import (
	"time"
)

// JobTime runs a Job and times it.
//
// Parameters:
//
//	job (Job): The job code that should be run.
//
// Returns:
//
//	time.Time: The time recorded before the job was run.
//	time.Time: The time recorded after the job was run.
func JobTime(job Job) (time.Time, time.Time) {
	startTime := time.Now()

	job()

	endTime := time.Now()

	return startTime, endTime
}
