package main

import (
	"time"

	"github.com/orsinium-labs/gamepad"
	"github.com/orsinium-labs/tellowerk/controllers"
	"go.uber.org/zap"
)

type GamePad struct {
	controller controllers.Controller
	gamepad    *gamepad.GamePad
	logger     *zap.Logger
	finish     chan<- struct{}
}

func (g *GamePad) Start() {
	go g.worker()
}

func (g *GamePad) worker() {
	var oldState gamepad.State
	time.Sleep(time.Second)
	g.logger.Debug("waiting input from gamepad")
	for {
		newState, err := g.gamepad.State()
		if err != nil {
			g.logger.Error("read gamepad", zap.Error(err))
			time.Sleep(2 * time.Second)
			continue
		}
		if newState.Start() {
			g.logger.Debug("closing connection")
			g.finish <- struct{}{}
			return
		}
		err = g.update(oldState, newState)
		if err != nil {
			g.logger.Error("execute gamepad command", zap.Error(err))
			time.Sleep(2 * time.Second)
			continue
		}
		oldState = newState
		time.Sleep(100 * time.Millisecond)
	}
}

func (g *GamePad) update(oldS, newS gamepad.State) error {
	var err error

	// nahdle take off and land
	if !oldS.A() && newS.A() {
		return g.controller.TakeOff()
	}
	if !oldS.B() && newS.B() {
		return g.controller.Land()
	}

	// handle ox rotation
	if oldS.LS().X != newS.LS().X {
		if newS.LS().X >= 0 {
			err = g.controller.Clockwise(newS.LS().X)
		} else {
			err = g.controller.CounterClockwise(-newS.LS().X)
		}
	}
	if err != nil {
		return err
	}

	// handle oy movement
	if oldS.LS().Y != newS.LS().Y {
		if newS.LS().Y >= 0 {
			err = g.controller.Backward(newS.LS().Y)
		} else {
			err = g.controller.Forward(-newS.LS().Y)
		}
	}
	if err != nil {
		return err
	}

	// handle ox movement
	if oldS.RS().X != newS.RS().X {
		if newS.RS().X >= 0 {
			err = g.controller.Right(newS.RS().X)
		} else {
			err = g.controller.Left(-newS.RS().X)
		}
	}
	if err != nil {
		return err
	}

	// handle oz movement
	if oldS.RS().Y != newS.RS().Y {
		if newS.RS().Y >= 0 {
			err = g.controller.Down(newS.RS().Y)
		} else {
			err = g.controller.Up(-newS.RS().Y)
		}
	}
	if err != nil {
		return err
	}

	return nil
}
