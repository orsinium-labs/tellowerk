package controllers

import "go.uber.org/zap"

type Logger struct {
	logger *zap.Logger
}

var _ Controller = &Logger{}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

func (c *Logger) Name() string {
	return "logger"
}

func (c *Logger) Start() error {
	c.logger.Debug("start")
	return nil
}

func (c *Logger) Stop() error {
	c.logger.Debug("stop")
	return nil
}

// Commands

func (c *Logger) TakeOff() error {
	c.logger.Debug("take off")
	return nil
}

func (c *Logger) Land() error {
	c.logger.Debug("land")
	return nil
}

func (c *Logger) Clockwise(val int) error {
	c.logger.Debug("rotate clockwise", zap.Int("val", val))
	return nil
}

func (c *Logger) CounterClockwise(val int) error {
	c.logger.Debug("rotate counter clockwise", zap.Int("val", val))
	return nil
}
