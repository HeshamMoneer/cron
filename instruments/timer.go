package instruments

import (
	"time"
)

func JobTime(job Job) (time.Time, time.Time) {
	startTime := time.Now()

	job()

	endTime := time.Now()

	return startTime, endTime
}
