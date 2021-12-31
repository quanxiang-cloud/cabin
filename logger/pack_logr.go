package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	//"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var _ AdaptedLogger = (*packLogr)(nil)

// ZapField return a zap field key-value pair
func ZapField(field string, val interface{}) zap.Field {
	return zapIt(field, val)
}

// packLogr adapt logr.Loger format like Logger
type packLogr struct {
	l               *zap.Logger
	minLevel        Level
	numericLevelKey string
	allowZapFields  bool
	panicMessages   bool
	keyAndValues    []zap.Field
}

func newPackLogr(cfg *Config, options ...Option) *packLogr {
	opts := append([]Option{
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zap.DPanicLevel)}, options...)
	zapLogger, err := zap.NewDevelopment(opts...)
	if err != nil {
		panic(err)
	}
	p := &packLogr{
		l:              zapLogger,
		minLevel:       _defaultLevel,
		allowZapFields: true,
		keyAndValues:   nil,
	}
	return p.init(cfg)
}

// init reset the log level of logger
func (p *packLogr) init(cfg *Config) *packLogr {
	lv := _defaultLevel.Int()

	env, err := GetLogLevelFromEnv()
	switch {
	case err == nil:
		lv = env
	case cfg != nil:
		lv = cfg.Level
	}

	return p.initLevel(lv)
}

func (p *packLogr) initLevel(lv int) *packLogr {
	if err := validateLogLevel(lv); err == nil {
		p.minLevel = Level(lv)
		p.l.WithOptions(zap.AddCallerSkip(2)).Info("init-logger-level", ZapField("level", lv))
	} else {
		p.PutError(err, "init-logger-level", zap.Int("level", lv))
	}

	return p
}

const noLevel = -1

// NOTE: from github.com/go-logr/zapr.Logger
//
// handleFields converts a bunch of arbitrary key-value pairs into Zap fields.  It takes
// additional pre-converted Zap fields, for use with automatically attached fields, like
// `error`.
func (p *packLogr) handleFields(lvl int, args []interface{}, additional ...zap.Field) []zap.Field {
	injectNumericLevel := p.numericLevelKey != "" && lvl != noLevel

	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields and no "v" field.
		if !injectNumericLevel {
			return additional
		}
		// Slightly slower fast path when we need to inject "v".
		return append(additional, zap.Int(p.numericLevelKey, lvl))
	}

	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	numFields := len(args)/2 + len(additional)
	if injectNumericLevel {
		numFields++
	}
	fields := make([]zap.Field, 0, numFields)
	if injectNumericLevel {
		fields = append(fields, zap.Int(p.numericLevelKey, lvl))
	}
	for i := 0; i < len(args); {
		// Check just in case for strongly-typed Zap fields,
		// which might be illegal (since it breaks
		// implementation agnosticism). If disabled, we can
		// give a better error message.
		if field, ok := args[i].(zap.Field); ok {
			if p.allowZapFields {
				fields = append(fields, field)
				i++
				continue
			}
			if p.panicMessages {
				p.l.WithOptions(zap.AddCallerSkip(1)).DPanic("strongly-typed Zap Field passed to logr", zapIt("zap field", args[i]))
			}
			break
		}

		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			if p.panicMessages {
				p.l.WithOptions(zap.AddCallerSkip(1)).DPanic("odd number of arguments passed as key-value pairs for logging", zapIt("ignored key", args[i]))
			}
			break
		}

		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			if p.panicMessages {
				p.l.WithOptions(zap.AddCallerSkip(1)).DPanic("non-string key argument passed to logging, ignoring all later arguments", zapIt("invalid key", key))
			}
			break
		}

		fields = append(fields, zapIt(keyStr, val))
		i += 2
	}

	return append(additional, fields...)
}

func zapIt(field string, val interface{}) zap.Field {
	// Handle types that implement logr.Marshaler: log the replacement
	// object instead of the original one.
	if marshaler, ok := val.(logr.Marshaler); ok {
		val = marshaler.MarshalLog()
	}
	return zap.Any(field, val)
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

	args := p.handleFields(lv.Int(), context, p.keyAndValues...)
	if cap(p.keyAndValues) < len(args) {
		p.keyAndValues = args[:len(p.keyAndValues)]
	}

	switch lv {
	case DebugLevel:
		p.l.Debug(msg, args...)
	case InfoLevel:
		p.l.Info(msg, args...)
	case WarnLevel:
		p.l.Warn(msg, args...)
	case ErrorLevel:
		p.l.Error(msg, args...)
	case DPanicLevel:
		p.l.DPanic(msg, args...)
	case PanicLevel:
		p.l.Panic(msg, args...)
	case FatalLevel:
		p.l.Fatal(msg, args...)
	}
}

func (p *packLogr) clone() *packLogr {
	return &packLogr{
		l:              p.l,
		minLevel:       p.minLevel,
		allowZapFields: p.allowZapFields,
		keyAndValues:   append([]zap.Field(nil), p.keyAndValues...),
	}
}

// WithValues returns a new Logger with additional key/value pairs.
func (p *packLogr) WithValues(keyAndValues ...interface{}) AdaptedLogger {
	n := p.clone()
	n.l = p.l.With(p.handleFields(noLevel, keyAndValues)...)
	return n
}

// WithName returns a new Logger with the specified name appended.
func (p *packLogr) WithName(name string) AdaptedLogger {
	n := p.clone()
	n.l = p.l.Named(name)
	return n
}

// WithLevel returns a new Logger with the specified level filter.
func (p *packLogr) WithLevel(level Level) AdaptedLogger {
	return p.clone().initLevel(level.Int())
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (p *packLogr) WithOptions(opts ...Option) AdaptedLogger {
	n := p.clone()
	n.l = p.l.WithOptions(opts...)
	return n
}

// PutError write log with error
func (p *packLogr) PutError(err error, msg string, keyAndValues ...interface{}) {
	if err != nil {
		keyAndValues = append([]interface{}{ZapField("error", err)}, keyAndValues...)
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
	return p.l.Sync()
}
