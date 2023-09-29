package instruments

import (
	"cron/scheduler"
)

func handle(job scheduler.Job, log logger, identifier int) scheduler.Job {
	wrappedJob := func() {
		defer func() {
			if err := recover(); err != nil {
				log.error(err)
				log.warn("Exited job", identifier)
			}
		}()

		job()
	}

	return wrappedJob
}
