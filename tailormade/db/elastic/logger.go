package elastic

import (
	"fmt"

	"github.com/go-logr/logr"
)

type logger struct {
	log logr.Logger
}

func newLogger(log logr.Logger) logger {
	return logger{
		log: log.WithName("elastic"),
	}
}

func (l logger) Printf(format string, v ...interface{}) {
	l.log.Info(fmt.Sprintf(format, v...))
}
