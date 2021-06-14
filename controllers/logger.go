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

func (c *Logger) Clockwise() error {
	c.logger.Debug("rotate clockwise")
	return nil
}

func (c *Logger) Land() error {
	c.logger.Debug("land")
	return nil
}
