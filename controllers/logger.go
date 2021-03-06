package controllers

import (
	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

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

// take off and land

func (c *Logger) TakeOff() error {
	c.logger.Debug("take off")
	return nil
}

func (c *Logger) ThrowTakeOff() error {
	c.logger.Debug("throw take off")
	return nil
}

func (c *Logger) Land() error {
	c.logger.Debug("land")
	return nil
}

func (c *Logger) PalmLand() error {
	c.logger.Debug("palm land")
	return nil
}

func (c *Logger) StopLanding() error {
	c.logger.Debug("stop landing")
	return nil
}

// movment

func (c *Logger) Rotate(val int) error {
	c.logger.Debug("rotate", zap.Int("val", val))
	return nil
}

func (c *Logger) OX(val int) error {
	c.logger.Debug("move ox", zap.Int("val", val))
	return nil
}

func (c *Logger) OY(val int) error {
	c.logger.Debug("move oy", zap.Int("val", val))
	return nil
}

func (c *Logger) OZ(val int) error {
	c.logger.Debug("move oz", zap.Int("val", val))
	return nil
}

func (c *Logger) Hover() error {
	// c.logger.Debug("hover")
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

func (c *Logger) Bounce() error {
	c.logger.Debug("bounce")
	return nil
}

func (c *Logger) SetFastMode() error {
	c.logger.Debug("fast mode")
	return nil
}

func (c *Logger) SetSlowMode() error {
	c.logger.Debug("slow mode")
	return nil
}

func (c *Logger) SetExposure(val int) error {
	c.logger.Debug("set exposure", zap.Int("val", val))
	return nil
}

func (c *Logger) SetVideoBitRate(val tello.VideoBitRate) error {
	c.logger.Debug("set video bitrate", zap.Int("val", int(val)))
	return nil
}
