package logging

import (
	"fmt"
	"log"
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
)

type Logger struct {
	name          string
	infoLogger    *log.Logger
	successLogger *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func NewLogger(name string) *Logger {
	prefix := fmt.Sprintf("[%s] ", name)
	return &Logger{
		name:          name,
		infoLogger:    log.New(os.Stdout, prefix+ColorBlue+"INFO: "+ColorReset, log.Ldate|log.Ltime),
		successLogger: log.New(os.Stdout, prefix+ColorGreen+"SUCCESS: "+ColorReset, log.Ldate|log.Ltime),
		warningLogger: log.New(os.Stdout, prefix+ColorYellow+"WARNING: "+ColorReset, log.Ldate|log.Ltime),
		errorLogger:   log.New(os.Stderr, prefix+ColorRed+"ERROR: "+ColorReset, log.Ldate|log.Ltime),
	}
}

func (l *Logger) Info(v ...interface{}) {
	_ = l.infoLogger.Output(2, fmt.Sprintln(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	_ = l.infoLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Warning(v ...interface{}) {
	_ = l.warningLogger.Output(2, fmt.Sprintln(v...))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	_ = l.warningLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	_ = l.errorLogger.Output(2, fmt.Sprintln(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	_ = l.errorLogger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Success(v ...interface{}) {
	_ = l.successLogger.Output(2, fmt.Sprintln(v...))
}

func (l *Logger) Successf(format string, v ...interface{}) {
	_ = l.successLogger.Output(2, fmt.Sprintf(format, v...))
}
