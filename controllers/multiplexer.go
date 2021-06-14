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
			return fmt.Errorf("start (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Stop() error {
	for _, sub := range c.controllers {
		err := sub.Stop()
		if err != nil {
			return fmt.Errorf("stop (%s): %v", c.Name(), err)
		}
	}
	return nil
}

// Commands

func (c *Multiplexer) TakeOff() error {
	for _, sub := range c.controllers {
		err := sub.TakeOff()
		if err != nil {
			return fmt.Errorf("take off (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Land() error {
	for _, sub := range c.controllers {
		err := sub.Land()
		if err != nil {
			return fmt.Errorf("land (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Clockwise(val int) error {
	for _, sub := range c.controllers {
		err := sub.Clockwise(val)
		if err != nil {
			return fmt.Errorf("rotate clockwise (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) CounterClockwise(val int) error {
	for _, sub := range c.controllers {
		err := sub.CounterClockwise(val)
		if err != nil {
			return fmt.Errorf("rotate counter clockwise (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Forward(val int) error {
	for _, sub := range c.controllers {
		err := sub.Forward(val)
		if err != nil {
			return fmt.Errorf("forward (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Backward(val int) error {
	for _, sub := range c.controllers {
		err := sub.Backward(val)
		if err != nil {
			return fmt.Errorf("backward (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Left(val int) error {
	for _, sub := range c.controllers {
		err := sub.Left(val)
		if err != nil {
			return fmt.Errorf("left (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Right(val int) error {
	for _, sub := range c.controllers {
		err := sub.Right(val)
		if err != nil {
			return fmt.Errorf("right (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Up(val int) error {
	for _, sub := range c.controllers {
		err := sub.Up(val)
		if err != nil {
			return fmt.Errorf("up (%s): %v", c.Name(), err)
		}
	}
	return nil
}

func (c *Multiplexer) Down(val int) error {
	for _, sub := range c.controllers {
		err := sub.Down(val)
		if err != nil {
			return fmt.Errorf("down (%s): %v", c.Name(), err)
		}
	}
	return nil
}
