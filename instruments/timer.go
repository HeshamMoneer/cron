package instruments

import (
	"cron/scheduler"
	"time"
)

func jobTime(job scheduler.Job) (time.Time, time.Time) {
	startTime := time.Now()

	job()

	endTime := time.Now()

	return startTime, endTime
}
