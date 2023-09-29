package scheduler

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	isConsole   bool
}

func getLoggerOutput(file ...string) io.Writer {
	if len(file) == 0 {
		return os.Stdout
	} else {
		fileName := file[0]
		file, _ := os.Create("./logs/" + fileName + ".log")
		return file
	}
}

// NewLogger creates a new logger with three logging levels Info, Warn, and Error.
//
// Parameters:
//
//	fileName [optional] (string): If passed, the logger creates a file in the logs folder and appends logs to it. Otherwise, the logger appends logs by default to the standard output.
func NewLogger(file ...string) Logger {
	output := getLoggerOutput(file...)
	flags := log.Ldate | log.Ltime | log.Lmicroseconds
	return Logger{
		infoLogger:  log.New(output, "CRON INFO : ", flags),
		warnLogger:  log.New(output, "CRON WARN : ", flags),
		errorLogger: log.New(output, "CRON ERROR: ", flags),
		isConsole:   output == os.Stdout,
	}
}

func (l *Logger) Info(v ...any) {
	l.infoLogger.Println(v...)
}

func (l *Logger) Warn(v ...any) {
	if l.isConsole {
		fmt.Printf("\033[1;33m")
		defer fmt.Printf("\033[0m")
	}
	l.warnLogger.Println(v...)
}

func (l *Logger) Error(v ...any) {
	if l.isConsole {
		fmt.Printf("\033[1;31m")
		defer fmt.Printf("\033[0m")
	}
	l.errorLogger.Println(v...)
}
