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

func (c *Driver) Hover() error {
	c.Driver.Hover()
	return nil
}

func (c *Driver) OX(val int) error {
	return c.Driver.Right(val)
}

func (c *Driver) OY(val int) error {
	return c.Driver.Forward(val)
}

func (c *Driver) OZ(val int) error {
	return c.Driver.Up(val)
}

func (c *Driver) Rotate(val int) error {
	return c.Driver.Clockwise(val)
}

func (c *Driver) SetVideoBitRate(val tello.VideoBitRate) error {
	return c.Driver.SetVideoEncoderRate(val)
}
