package logger

import (
	"log"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() Logger {
	return Logger{logger: log.Default()}
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Printf("[INFO] "+format, args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.logger.Printf("[WARNING] "+format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Printf("[ERROR] "+format, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.Fatalf("[FATAL] "+format, args...)
}

func (l *Logger) Panicf(format string, args ...any) {
	l.logger.Panicf("[PANIC] "+format, args...)
}
