package instruments

func HandleErrors(job Job, log Logger, identifier int) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			log.Warn("Exited job", identifier)
		}
	}()

	job()
}
