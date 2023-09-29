package scheduler

// HandleErrors runs a Job while guarding the program execution from any errors in the job.
//
// Parameters:
//	job (Job): The job code that should be run.
//	identifier (int): The cron job identifier.
//	log (Logger): The logger of the program.
func HandleErrors(job Job, identifier int, log Logger) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			log.Warn("Exited job", identifier)
		}
	}()

	job()
}
