package elastic

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/olivere/elastic/v7"
)

// Config config
type Config struct {
	Host []string
	Log  bool
}

// NewClient new elasticsearch client
func NewClient(conf *Config, log logr.Logger) (*elastic.Client, error) {
	opts := make([]elastic.ClientOptionFunc, 0, len(conf.Host))
	for _, host := range conf.Host {
		opts = append(opts, elastic.SetURL(host))
	}
	opts = append(opts, elastic.SetErrorLog(newLogger(log)))
	if conf.Log {
		opts = append(opts, elastic.SetInfoLog(newLogger(log)))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	_, _, err = client.Ping(conf.Host[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = client.ElasticsearchVersion(conf.Host[0])
	if err != nil {
		return nil, err
	}

	return client, nil
}
