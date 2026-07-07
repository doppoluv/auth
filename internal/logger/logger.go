package logger

import (
	"log"

	"auth/internal/logger/interfaces"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() interfaces.Logger {
	return &Logger{logger: log.Default()}
}

func (l *Logger) Printf(format string, args ...any) {
	l.logger.Printf("[INFO] "+format, args...)
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
