package controllers

import "fmt"

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

func (c *Multiplexer) Clockwise(val int) error {
	for _, sub := range c.controllers {
		err := sub.Clockwise(val)
		if err != nil {
			return fmt.Errorf("rotate clockwise (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) CounterClockwise(val int) error {
	for _, sub := range c.controllers {
		err := sub.CounterClockwise(val)
		if err != nil {
			return fmt.Errorf("rotate counter clockwise (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Forward(val int) error {
	for _, sub := range c.controllers {
		err := sub.Forward(val)
		if err != nil {
			return fmt.Errorf("forward (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Backward(val int) error {
	for _, sub := range c.controllers {
		err := sub.Backward(val)
		if err != nil {
			return fmt.Errorf("backward (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Left(val int) error {
	for _, sub := range c.controllers {
		err := sub.Left(val)
		if err != nil {
			return fmt.Errorf("left (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Right(val int) error {
	for _, sub := range c.controllers {
		err := sub.Right(val)
		if err != nil {
			return fmt.Errorf("right (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Up(val int) error {
	for _, sub := range c.controllers {
		err := sub.Up(val)
		if err != nil {
			return fmt.Errorf("up (%s): %v", sub.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Down(val int) error {
	for _, sub := range c.controllers {
		err := sub.Down(val)
		if err != nil {
			return fmt.Errorf("down (%s): %v", sub.Name(), err)
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
