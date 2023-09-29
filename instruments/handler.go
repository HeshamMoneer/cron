package instruments

func handle(job Job, log logger, identifier int) Job {
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
