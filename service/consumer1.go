package service

import "go.uber.org/zap"

type Consumer1 struct {
	log *zap.Logger
}

func NewConsumer1(log *zap.Logger) *Consumer1 {
	return &Consumer1{
		log: log,
	}
}

func (c *Consumer1) Consume(msg string) error {
	c.log.Info("Incoming message ", zap.String("a", msg))
	return nil
}
