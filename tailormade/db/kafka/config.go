package kafka

import (
	"crypto/tls"

	"github.com/Shopify/sarama"
)

// Config config
type Config struct {
	Sarama sarama.Config `yaml:"sarama"`

	Broker []string    `yaml:"broker"`
	TLS    *tls.Config `yaml:"tls"`
}

func pre(conf Config) *sarama.Config {
	config := sarama.NewConfig()

	// TLS
	config.Net.TLS.Enable = conf.Sarama.Net.TLS.Enable
	config.Net.TLS.Config = conf.Sarama.Net.TLS.Config

	return config
}
