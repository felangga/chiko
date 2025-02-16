package logger

import "github.com/felangga/chiko/internal/entity"

type Logger struct {
	logChannel    chan entity.Log
	outputChannel chan entity.Output
}

func New() *Logger {
	return &Logger{
		logChannel:    make(chan entity.Log, 100),
		outputChannel: make(chan entity.Output, 100),
	}
}

func (l *Logger) Channel() chan entity.Log {
	return l.logChannel
}

func (l *Logger) OutputChannel() chan entity.Output {
	return l.outputChannel
}

func (l *Logger) Info(message string) {
	log := entity.Log{
		Content: message,
		Type:    entity.LOG_INFO,
	}
	l.logChannel <- log
}

func (l *Logger) Warning(message string) {
	log := entity.Log{
		Content: message,
		Type:    entity.LOG_WARNING,
	}
	l.logChannel <- log
}

func (l *Logger) Error(message string) {
	log := entity.Log{
		Content: message,
		Type:    entity.LOG_ERROR,
	}
	l.logChannel <- log
}

func (l *Logger) PrintOutput(output entity.Output) {
	l.outputChannel <- output
}
