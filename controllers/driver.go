package controllers

import "gobot.io/x/gobot/platforms/dji/tello"

type Driver struct {
	*tello.Driver
}

var _ Controller = &Driver{}

func NewDriver(driver *tello.Driver) *Driver {
	return &Driver{Driver: driver}
}

func (c *Driver) Name() string {
	return "driver"
}

func (c *Driver) Stop() error {
	return c.Driver.Halt()
}
