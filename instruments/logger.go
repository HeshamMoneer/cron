package instruments

import (
	"io"
	"log"
	"os"
)

type logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func getLoggerOutput(args ...string) io.Writer {
	if len(args) == 0 {
		return os.Stdout
	} else {
		fileName := args[0]
		file, _ := os.Create("./logs/" + fileName + ".log")
		return file
	}
}

func Logger(args ...string) logger {
	output := getLoggerOutput(args...)
	flags := log.Ldate | log.Ltime | log.Lmicroseconds
	return logger{
		infoLogger:  log.New(output, "CRON INFO : ", flags),
		warnLogger:  log.New(output, "CRON WARN : ", flags),
		errorLogger: log.New(output, "CRON ERROR: ", flags),
	}
}

func (l *logger) info(v ...any) {
	l.infoLogger.Println(v...)
}

func (l *logger) warn(v ...any) {
	l.warnLogger.Println(v...)
}

func (l *logger) error(v ...any) {
	l.errorLogger.Println(v...)
}
