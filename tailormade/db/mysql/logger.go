/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysql

import (
	"context"
	"fmt"
	"time"

	pkglogger "github.com/quanxiang-cloud/cabin/logger"

	lg "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// logger gorm logger
type logger struct {
	log pkglogger.AdaptedLogger
	lg.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// newLogger new gorm logger
func newLogger(log pkglogger.AdaptedLogger, config lg.Config) lg.Interface {
	log = log.WithName("gorm")
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &logger{
		log:          log,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *logger) LogMode(level lg.LogLevel) lg.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func msssage(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Info {
		l.log.Info("SQL", msssage(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Warn {
		l.log.Info("warn", "SQL", msssage(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Error {
		l.log.Info("error", "SQL", msssage(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= lg.Error:
			sql, rows := fc()
			if rows == -1 {
				l.log.Info("debug", "SQL", msssage(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.log.Info("debug", "SQL", msssage(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= lg.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.log.Info("debug", "SQL", msssage(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.log.Info("debug", "SQL", msssage(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		case l.LogLevel >= lg.Info:
			sql, rows := fc()
			if rows == -1 {
				l.log.Info("debug", "SQL", msssage(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.log.Info("debug", "SQL", msssage(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		}
	}
}
