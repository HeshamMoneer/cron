package instruments

import (
	"time"
)

func jobTime(job Job) (time.Time, time.Time) {
	startTime := time.Now()

	job()

	endTime := time.Now()

	return startTime, endTime
}
