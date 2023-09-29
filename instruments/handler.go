package instruments

func HandleErrors(job Job, identifier int, log Logger) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			log.Warn("Exited job", identifier)
		}
	}()

	job()
}
