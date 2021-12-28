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
