package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

// New создаёт zap-логгер (режим development, чтобы видеть читаемые логи).
func New() *Logger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	return &Logger{log}
}

// Sync выполняет корректное завершение логгера.
func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}
