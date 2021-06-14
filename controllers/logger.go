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

func (c *Logger) Forward(val int) error {
	c.logger.Debug("forward", zap.Int("val", val))
	return nil
}

func (c *Logger) Backward(val int) error {
	c.logger.Debug("backward", zap.Int("val", val))
	return nil
}

func (c *Logger) Left(val int) error {
	c.logger.Debug("left", zap.Int("val", val))
	return nil
}

func (c *Logger) Right(val int) error {
	c.logger.Debug("right", zap.Int("val", val))
	return nil
}

func (c *Logger) Up(val int) error {
	c.logger.Debug("up", zap.Int("val", val))
	return nil
}

func (c *Logger) Down(val int) error {
	c.logger.Debug("down", zap.Int("val", val))
	return nil
}

func (c *Logger) LeftFlip() error {
	c.logger.Debug("left flip")
	return nil
}

func (c *Logger) RightFlip() error {
	c.logger.Debug("right flip")
	return nil
}

func (c *Logger) FrontFlip() error {
	c.logger.Debug("front flip")
	return nil
}

func (c *Logger) BackFlip() error {
	c.logger.Debug("back flip")
	return nil
}
