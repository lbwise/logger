package logger

import (
	"errors"
	"io"
	"os"
	"time"
)

var (
	CannotWriteError   = errors.New("cannot write log to output")
	InvalidOutputError = errors.New("nil output writer provided")
)

type Logger struct {
	logs        []Log
	out         io.Writer
	colored     bool
	prefix      string
	includeTime bool
	time        func() time.Time
}

func New(options ...Option) *Logger {
	// Create logger with default cfg
	logger := &Logger{
		out:         os.Stdout,
		logs:        make([]Log, 0),
		prefix:      "",
		includeTime: false,
		colored:     false,
		time:        time.Now,
	}

	// Apply options
	for _, opt := range options {
		opt(logger)
	}
	return logger
}

func (l *Logger) Log(s string) error {
	return l.addLog(InfoLog(s))
}

func (l *Logger) Warn(message string, severity WarnSeverity) error {
	return l.addLog(WarnLog{severity, message})
}

func (l *Logger) Error(err error, code int) error {
	return l.addLog(ErrorLog{err, code})
}

func (l *Logger) addLog(log Log) error {
	if l.out == nil {
		return InvalidOutputError
	}

	var msg string
	if l.prefix != "" && l.includeTime {
		msg += l.prefix + "/" + l.time().String()[:19] + " "
	} else if l.prefix != "" {
		msg += l.prefix + " "
	} else if l.includeTime && l.prefix == "" {
		msg += l.time().String()[:19] + " "
	}

	msg += log.Write(l.colored)

	if _, err := l.out.Write([]byte(msg)); err != nil {
		return CannotWriteError
	}
	l.logs = append(l.logs, log)
	return nil
}

type Option func(*Logger)

func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.out = w
	}
}

func WithPrefix(prefix string) Option {
	return func(l *Logger) {
		l.prefix = prefix
	}
}

func WithTimeIncluded() Option {
	return func(l *Logger) {
		l.includeTime = true
	}
}

func WithColor() Option {
	return func(l *Logger) {
		l.colored = true
	}
}

func WithClock(c func() time.Time) Option {
	return func(l *Logger) {
		l.time = c
	}
}
