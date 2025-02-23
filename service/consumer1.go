package service

import (
	"go.uber.org/zap"
	"uber-fx-workshop/config"
)

type Consumer1 struct {
	log    *zap.Logger
	config *config.RootConfig
}

func NewConsumer1(log *zap.Logger, config *config.RootConfig) *Consumer1 {
	return &Consumer1{
		log:    log,
		config: config,
	}
}

func (c *Consumer1) Consume(msg string) error {
	c.log.Info("Incoming message ", zap.String("a", msg))
	for consumerName, config := range c.config.KafkaAthena.Consumers {
		c.log.Info("config", zap.String(consumerName, config.Cluster))
		c.log.Info("config", zap.String(consumerName, config.Topic))
	}
	return nil
}
