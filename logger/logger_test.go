package logger_test

import (
	"errors"
	"os"
	"testing"

	"github.com/quanxiang-cloud/cabin/logger"
)

func TestLogger(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Warn("TestLogger recovered")
		}
	}()

	os.Setenv(logger.EnvLogLevel, "0")

	cfg := &logger.Config{
		Level: logger.DebugLevel.Int(),
	}
	logger.Logger = logger.New(cfg)

	namedLog := logger.Logger.WithName("named")
	log := logger.Logger
	valueLog := logger.Logger.WithValues(
		logger.ZapField("val", "foo"),
		logger.ZapField("val2", "bar"),
	)

	log.Sync()
	namedLog.Infof("info %s", "foo")
	namedLog.PutError(nil, "without-error")
	namedLog.PutError(errors.New("errMessage"), "with-error")
	log.Infof("info %s", "foo")
	log.Debugf("debug")
	log.Warnf("warn")
	log.Errorf("err")
	log.Debug("debug")
	log.Warn("warn")
	log.Error("err")
	namedLog.Debugw("debug", "foo", 1)
	namedLog.Infow("info", "foo", 1)
	namedLog.Warnw("warn", logger.ZapField("foo", "zaped"))
	valueLog.Debugw("debug", "foo", 1)
	valueLog.Infow("info", "foo", 1)
	valueLog.Warnw("warn", logger.ZapField("foo", "zaped"))
	log.Panic("panic")
	log.Info("info2")
	log.Fatal("fatal")
}

func TestPanic(t *testing.T) {
	test := func(lv logger.Level, method string) {
		defer func() { recover() }()
		switch lv {
		case logger.DPanicLevel:
			switch method {
			case "f":
				logger.Logger.DPanicf("DPanic")
			case "w":
				logger.Logger.DPanicw("DPanic")
			case "":
				logger.Logger.DPanic("DPanic")
			}

		case logger.PanicLevel:
			switch method {
			case "f":
				logger.Logger.Panicf("Panic")
			case "w":
				logger.Logger.Panicw("Panic")
			case "":
				logger.Logger.Panic("Panic")
			}
		}
	}
	test(logger.DPanicLevel, "f")
	test(logger.DPanicLevel, "w")
	test(logger.DPanicLevel, "")
	test(logger.PanicLevel, "f")
	test(logger.PanicLevel, "w")
	test(logger.PanicLevel, "")
}
