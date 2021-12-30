package logger

import (
	"fmt"
	"os"
	"strconv"
)

// DefaultLogLevel is the *number* that control the level of Logger with highest priority.
//
// NOTE: ***KEEP EMPTY VALUE HERE***
// This value is expected to modify by:
// GO run -ldflags "-X 'github.com/quanxiang-cloud/cabin/logger.DefaultLogLevel=0'"
const DefaultLogLevel = ""

// EnvLogLevel is the environment variable to control logger.LogLevel
// eg: "set CABIN_LOG_LEVEL=0"
const EnvLogLevel = "CABIN_LOG_LEVEL"

// GetLogLevelFromEnv get cabin log level from os.Getenv
func GetLogLevelFromEnv() (int, error) {
	env := DefaultLogLevel // default log level value
	if DefaultLogLevel == "" {
		env = os.Getenv(EnvLogLevel)
	}

	i, err := int64(0), error(nil)
	if i, err = strconv.ParseInt(env, 10, 8); err == nil {
		err = validateLogLevel(int(i))
	}

	if err != nil {
		if env != "" {
			fmt.Printf("**invalid int8 value %q of DefaultLogLevel or env.CABIN_LOG_LEVEL : %v\n", env, err)
		}
		return _defaultLevel.Int(), err
	}

	return int(i), nil
}
