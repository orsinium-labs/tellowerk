package controllers

import (
	"fmt"

	"gobot.io/x/gobot/platforms/dji/tello"
)

type Multiplexer struct {
	controllers []Controller
}

var _ Controller = &Multiplexer{}

func NewMultiplexer() *Multiplexer {
	return &Multiplexer{controllers: make([]Controller, 0)}
}

func (c *Multiplexer) Add(sub Controller) {
	c.controllers = append(c.controllers, sub)
}

func (c *Multiplexer) Name() string {
	return "multiplexer"
}

func (c *Multiplexer) Start() error {
	for _, sub := range c.controllers {
		err := sub.Start()
		if err != nil {
			return fmt.Errorf("start (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Stop() error {
	for _, sub := range c.controllers {
		err := sub.Stop()
		if err != nil {
			return fmt.Errorf("stop (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

// Commands

func (c *Multiplexer) TakeOff() error {
	for _, sub := range c.controllers {
		err := sub.TakeOff()
		if err != nil {
			return fmt.Errorf("take off (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) ThrowTakeOff() error {
	for _, sub := range c.controllers {
		err := sub.ThrowTakeOff()
		if err != nil {
			return fmt.Errorf("throw take off (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Land() error {
	for _, sub := range c.controllers {
		err := sub.Land()
		if err != nil {
			return fmt.Errorf("land (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) PalmLand() error {
	for _, sub := range c.controllers {
		err := sub.PalmLand()
		if err != nil {
			return fmt.Errorf("palm land (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) StopLanding() error {
	for _, sub := range c.controllers {
		err := sub.StopLanding()
		if err != nil {
			return fmt.Errorf("stop landing (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

// movement

// Rotate clockwise or, if negate, counter-clockwise.
func (c *Multiplexer) Rotate(val int) error {
	for _, sub := range c.controllers {
		err := sub.Rotate(val)
		if err != nil {
			return fmt.Errorf("rotate clockwise (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

// Move forward or, if negative, backward.
func (c *Multiplexer) OY(val int) error {
	for _, sub := range c.controllers {
		err := sub.OY(val)
		if err != nil {
			return fmt.Errorf("forward (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) OX(val int) error {
	for _, sub := range c.controllers {
		err := sub.OX(val)
		if err != nil {
			return fmt.Errorf("right (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) OZ(val int) error {
	for _, sub := range c.controllers {
		err := sub.OZ(val)
		if err != nil {
			return fmt.Errorf("up (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Hover() error {
	for _, sub := range c.controllers {
		err := sub.Hover()
		if err != nil {
			return fmt.Errorf("hover (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) LeftFlip() error {
	for _, sub := range c.controllers {
		err := sub.LeftFlip()
		if err != nil {
			return fmt.Errorf("left flip (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) RightFlip() error {
	for _, sub := range c.controllers {
		err := sub.RightFlip()
		if err != nil {
			return fmt.Errorf("right flip (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) FrontFlip() error {
	for _, sub := range c.controllers {
		err := sub.FrontFlip()
		if err != nil {
			return fmt.Errorf("front flip (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) BackFlip() error {
	for _, sub := range c.controllers {
		err := sub.BackFlip()
		if err != nil {
			return fmt.Errorf("back flip (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Bounce() error {
	for _, sub := range c.controllers {
		err := sub.Bounce()
		if err != nil {
			return fmt.Errorf("bounce (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) SetFastMode() error {
	for _, sub := range c.controllers {
		err := sub.SetFastMode()
		if err != nil {
			return fmt.Errorf("fast mode (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) SetSlowMode() error {
	for _, sub := range c.controllers {
		err := sub.SetSlowMode()
		if err != nil {
			return fmt.Errorf("slow mode (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) SetExposure(val int) error {
	for _, sub := range c.controllers {
		err := sub.SetExposure(val)
		if err != nil {
			return fmt.Errorf("set exposure (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) SetVideoBitRate(val tello.VideoBitRate) error {
	for _, sub := range c.controllers {
		err := sub.SetVideoBitRate(val)
		if err != nil {
			return fmt.Errorf("set video bitrate (%s): %v", sub.Name(), err)
		}
	}
	return nil
}
