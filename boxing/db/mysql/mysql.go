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
	"fmt"
	"time"

	pkglogger "github.com/quanxiang-cloud/cabin/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

// Config mysql config
type Config struct {
	Host     string
	DB       string
	User     string
	Password string
	Log      bool

	MaxIdleConns int
	MaxOpenConns int

	dsn string
}

const (
	// DSN_DEFAULT utf-8
	DSN_DEFAULT = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"

	// DSN_UTF8MB4 utf-8 mb4
	DSN_UTF8MB4 = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

// String return dsn string
func (c *Config) String() string {
	if c.dsn == "" {
		c.dsn = DSN_DEFAULT
	}
	return c.string(c.dsn)
}

// SetDSN set dsn
func (c *Config) SetDSN(dsn string) {
	c.dsn = dsn
}

func (c *Config) string(format string) string {
	return fmt.Sprintf(format, c.User, c.Password, c.Host, c.DB)
}

// New return a gorm db
func New(config Config, log pkglogger.AdaptedLogger) (*gorm.DB, error) {
	var logger lg.Interface
	if config.Log {
		logger = newLogger(log, lg.Config{
			SlowThreshold: time.Second,
			LogLevel:      lg.Info,
			Colorful:      true,
		})
	}
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       config.String(),
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}),
		&gorm.Config{
			Logger: logger,
		})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 20
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	return db, nil
}
