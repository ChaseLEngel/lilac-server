package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	logfile string
}

func New(logfile string) *Logger {
	logger := new(Logger)
	logger.logfile = logfile
	logger.write(format("INFO", "Logger started"))
	return logger
}

func (l *Logger) write(text string) {
	file, err := os.OpenFile(l.logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Logger failed to open file %v: %v\n", l.logfile, err)
		return
	}
	defer file.Close()
	if _, err = file.WriteString(text); err != nil {
		fmt.Printf("Logger failed to write to file %v: %v\n", l.logfile, err)
	}
}

func format(level string, text string) string {
	t := time.Now()
	formattedTime := t.Format("Mon Jan _2 15:04:05 2006")
	return fmt.Sprintf("[%v] [%v] %v", level, formattedTime, text)
}

func (l *Logger) Info(v ...interface{}) {
	l.write(format("INFO", fmt.Sprintf("%v\n", v...)))
}

func (l *Logger) Infof(s string, v ...interface{}) {
	l.write(format("INFO", fmt.Sprintf(s, v...)))
}

func (l *Logger) Warn(v ...interface{}) {
	l.write(format("WARN", fmt.Sprintf("%v\n", v...)))
}

func (l *Logger) Warnf(s string, v ...interface{}) {
	l.write(format("WARN", fmt.Sprintf(s, v...)))
}

func (l *Logger) Error(v ...interface{}) {
	l.write(format("ERRO", fmt.Sprintf("%v\n", v...)))
}

func (l *Logger) Errorf(s string, v ...interface{}) {
	l.write(format("ERRO", fmt.Sprintf(s, v...)))
}
