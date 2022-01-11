package logger

import (
	"fmt"

	"github.com/go-logr/logr"
)

var _ AdaptedLogger = (*packLogr)(nil)

type packLogr struct {
	l        logr.Logger
	minLevel Level
}

func newPackLogr(log logr.Logger) *packLogr {
	p := &packLogr{
		l:        log,
		minLevel: noLevel,
	}
	return p
}

func (p *packLogr) log(lv Level, format string, fmtArgs []interface{}, context []interface{}) {
	if lv <= _maxIgnoreLevel && lv < p.minLevel { //ignore level with too low priority
		return
	}

	msg := format
	switch {
	case format == "" && len(fmtArgs) > 0:
		msg = fmt.Sprint(fmtArgs...)
	case format != "" && len(fmtArgs) > 0:
		msg = fmt.Sprintf(format, fmtArgs...)
	}

	if lv < ErrorLevel {
		p.l.Info(msg, context...)
	} else {
		p.l.Error(nil, msg, context...)
	}
}

func (p *packLogr) clone() *packLogr {
	return &packLogr{
		l:        p.l,
		minLevel: p.minLevel,
	}
}

// WithValues returns a new Logger with additional key/value pairs.
func (p *packLogr) WithValues(keyAndValues ...interface{}) AdaptedLogger {
	n := p.clone()
	n.l = p.l.WithValues(keyAndValues...)
	return n
}

// WithName returns a new Logger with the specified name appended.
func (p *packLogr) WithName(name string) AdaptedLogger {
	n := p.clone()
	n.l = p.l.WithName(name)
	return n
}

func (p *packLogr) initLevel(lv int) *packLogr {
	if err := validateLogLevel(lv); err == nil {
		p.minLevel = Level(lv)
		p.l.WithCallDepth(2).Info("init-logger-level", ZapField("level", lv))
	} else {
		p.PutError(err, "init-logger-level", ZapField("level", lv))
	}

	return p
}

// WithLevel returns a new Logger with the specified level filter.
func (p *packLogr) WithLevel(level Level) AdaptedLogger {
	return p.clone().initLevel(level.Int())
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (p *packLogr) WithOptions(opts ...Option) AdaptedLogger {
	n := p.clone()
	return n
}

// PutError write log with error
func (p *packLogr) PutError(err error, msg string, keyAndValues ...interface{}) {
	if err != nil {
		keyAndValues = append([]interface{}{"error", err}, keyAndValues...)
	}
	p.log(ErrorLevel, msg, nil, keyAndValues)
}

// Debug uses fmt.Sprint to construct and log a message.
func (p *packLogr) Debug(args ...interface{}) {
	p.log(DebugLevel, "", args, nil)
}

// Info uses fmt.Sprint to construct and log a message.
func (p *packLogr) Info(args ...interface{}) {
	p.log(InfoLevel, "", args, nil)
}

// Warn uses fmt.Sprint to construct and log a message.
func (p *packLogr) Warn(args ...interface{}) {
	p.log(WarnLevel, "", args, nil)
}

// Error uses fmt.Sprint to construct and log a message.
func (p *packLogr) Error(args ...interface{}) {
	p.log(ErrorLevel, "", args, nil)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (p *packLogr) DPanic(args ...interface{}) {
	p.log(DPanicLevel, "", args, nil)
}

// Panic uses fmt.Sprint to construct and log a message, then panicl.
func (p *packLogr) Panic(args ...interface{}) {
	p.log(PanicLevel, "", args, nil)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls ol.Exit.
func (p *packLogr) Fatal(args ...interface{}) {
	p.log(FatalLevel, "", args, nil)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (p *packLogr) Debugf(template string, args ...interface{}) {
	p.log(DebugLevel, template, args, nil)
}

// Infof uses fmt.Sprintf to log a templated message.
func (p *packLogr) Infof(template string, args ...interface{}) {
	p.log(InfoLevel, template, args, nil)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (p *packLogr) Warnf(template string, args ...interface{}) {
	p.log(WarnLevel, template, args, nil)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (p *packLogr) Errorf(template string, args ...interface{}) {
	p.log(ErrorLevel, template, args, nil)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (p *packLogr) DPanicf(template string, args ...interface{}) {
	p.log(DPanicLevel, template, args, nil)
}

// Panicf uses fmt.Sprintf to log a templated message, then panicl.
func (p *packLogr) Panicf(template string, args ...interface{}) {
	p.log(PanicLevel, template, args, nil)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls ol.Exit.
func (p *packLogr) Fatalf(template string, args ...interface{}) {
	p.log(FatalLevel, template, args, nil)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  l.With(keysAndValues).Debug(msg)
func (p *packLogr) Debugw(msg string, keysAndValues ...interface{}) {
	p.log(DebugLevel, msg, nil, keysAndValues)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (p *packLogr) Infow(msg string, keysAndValues ...interface{}) {
	p.log(InfoLevel, msg, nil, keysAndValues)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (p *packLogr) Warnw(msg string, keysAndValues ...interface{}) {
	p.log(WarnLevel, msg, nil, keysAndValues)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (p *packLogr) Errorw(msg string, keysAndValues ...interface{}) {
	p.log(ErrorLevel, msg, nil, keysAndValues)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (p *packLogr) DPanicw(msg string, keysAndValues ...interface{}) {
	p.log(DPanicLevel, msg, nil, keysAndValues)
}

// Panicw logs a message with some additional context, then panicl. The
// variadic key-value pairs are treated as they are in With.
func (p *packLogr) Panicw(msg string, keysAndValues ...interface{}) {
	p.log(PanicLevel, msg, nil, keysAndValues)
}

// Fatalw logs a message with some additional context, then calls ol.Exit. The
// variadic key-value pairs are treated as they are in With.
func (p *packLogr) Fatalw(msg string, keysAndValues ...interface{}) {
	p.log(FatalLevel, msg, nil, keysAndValues)
}

// Sync flushes any buffered log entries.
func (p *packLogr) Sync() error {
	return nil
}
