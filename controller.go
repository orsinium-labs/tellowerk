package main

import (
	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Controller struct {
	driver *tello.Driver
	logger *zap.Logger
}

func (c *Controller) Start() error {
	c.logger.Debug("start")
	return c.driver.Start()
}

func (c *Controller) Stop() error {
	c.logger.Debug("stop")
	return c.driver.Halt()
}

// Commands

func (c *Controller) TakeOff() error {
	c.logger.Debug("take off")
	return c.driver.TakeOff()
}

func (c *Controller) Land() error {
	c.logger.Debug("land")
	return c.driver.Land()
}
