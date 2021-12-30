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

	namedLog.Infof("info %s", "foo")
	namedLog.PutError(nil, "without-error")
	namedLog.PutError(errors.New("errMessage"), "with-error")
	log.Infof("info %s", "foo")
	log.Info("info")
	log.Debug("debug")
	log.Warn("warn")
	log.Error("err")
	namedLog.Debugw("debug", "foo", 1)
	namedLog.Infow("info", "foo", 1)
	namedLog.Warnw("warn", logger.ZapField("foo", "zaped"))
	log.Panic("panic")
	log.Info("info2")
	log.Fatal("fatal")
}
