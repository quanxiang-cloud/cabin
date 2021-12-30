package logger

import (
	"go.uber.org/zap/zapcore"
)

// A Level is a logging priority. Higher levels are more important.
type Level zapcore.Level

// Zap return zap level
func (l Level) Zap() zapcore.Level {
	return zapcore.Level(l)
}

// Int convert level to int
func (l Level) Int() int {
	return int(l)
}

func (l Level) index() int {
	return int(l - _minLevel)
}

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = Level(zapcore.DebugLevel)
	// InfoLevel is the default logging priority.
	InfoLevel = Level(zapcore.InfoLevel)
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = Level(zapcore.WarnLevel)
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = Level(zapcore.ErrorLevel)
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = Level(zapcore.DPanicLevel)
	// PanicLevel logs a message, then panics.
	PanicLevel = Level(zapcore.PanicLevel)
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = Level(zapcore.FatalLevel)

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)
