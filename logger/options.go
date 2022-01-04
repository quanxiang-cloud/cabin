package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// An Option configures a Logger.
type Option = zap.Option

// Field is an alias for Field. Aliasing this type dramatically
// improves the navigability of this package's API documentation.
type Field = zap.Field

// LevelEnabler decides whether a given logging level is enabled when logging a
// message.
//
// Enablers are intended to be used to implement deterministic filters;
// concerns like sampling are better implemented as a Core.
//
// Each concrete Level value implements a static LevelEnabler which returns
// true for itself and all higher logging levels. For example WarnLevel.Enabled()
// will return true for WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, and
// FatalLevel, but return false for InfoLevel and DebugLevel.
type LevelEnabler = zapcore.LevelEnabler

// Clock is a source of time for logged entries.
type Clock = zapcore.Clock

// Fields adds fields to the Logger.
func Fields(fs ...Field) Option {
	return zap.Fields(fs...)
}

// WithCaller configures the Logger to annotate each message with the filename,
// line number, and function name of zap's caller, or not, depending on the
// value of enabled. This is a generalized form of AddCaller.
func WithCaller(enabled bool) Option {
	return zap.WithCaller(enabled)
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option). When building wrappers around the
// Logger and SugaredLogger, supplying this Option prevents zap from always
// reporting the wrapper code as the caller.
func AddCallerSkip(skip int) Option {
	return zap.AddCallerSkip(skip)
}

// AddStacktrace configures the Logger to record a stack trace for all messages at
// or above a given level.
func AddStacktrace(lvl LevelEnabler) Option {
	return zap.AddStacktrace(lvl)
}

// WithClock specifies the clock used by the logger to determine the current
// time for logged entries. Defaults to the system clock with time.Now.
func WithClock(clock Clock) Option {
	return zap.WithClock(clock)
}
