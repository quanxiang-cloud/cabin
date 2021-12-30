package logger

import (
	"os"
	"strconv"
)

// EnvLogLevel is the environment variable to control cabin.logger.LogLevel
// eg: "set CABIN_LOG_LEVEL=0"
const EnvLogLevel = "CABIN_LOG_LEVEL"

// GetLogLevelFromEnv get cabin log level from os.Getenv
func GetLogLevelFromEnv() (int, error) {
	env := os.Getenv(EnvLogLevel)

	i, err := strconv.ParseInt(env, 10, 8)
	if err != nil {
		return DebugLevel.Int(), err
	}
	return int(i), nil
}
