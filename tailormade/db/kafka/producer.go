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
package kafka

import (
	"github.com/Shopify/sarama"
)

// NewSyncProducer new sync producer
func NewSyncProducer(conf Config) (sarama.SyncProducer, error) {
	config := pre(conf)
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.Broker, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// NewAsyncProducer new async producer
func NewAsyncProducer(conf Config) (sarama.AsyncProducer, error) {
	config := pre(conf)
	config.Producer.Return.Successes = conf.Sarama.Producer.Return.Successes

	producer, err := sarama.NewAsyncProducer(conf.Broker, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
