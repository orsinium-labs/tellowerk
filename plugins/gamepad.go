package plugins

import (
	"time"

	"github.com/orsinium-labs/gamepad"
	"github.com/orsinium-labs/tellowerk/controllers"
	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type GamePad struct {
	controller controllers.Controller
	info       *FlightInfo
	logger     *zap.Logger

	gamepad *gamepad.GamePad
	finish  chan struct{}
}

func NewGamePad(g *gamepad.GamePad) *GamePad {
	return &GamePad{
		gamepad: g,
		finish:  make(chan struct{}),
	}
}

func (g *GamePad) Connect(pl *Plugins) {
	g.controller = pl.Controller
	g.info = pl.FlightInfo
	g.logger = pl.Logger
}

func (g *GamePad) Start() error {
	go g.worker()
	return nil
}

func (GamePad) Stop() error {
	return nil
}

func (g *GamePad) Wait() {
	<-g.finish
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

	// handle take off and land
	if g.info.Flying() {
		if !oldS.A() && newS.A() {
			return g.controller.Land()
		}
		if oldS.A() && !newS.A() {
			return g.controller.StopLanding()
		}
		if !oldS.B() && newS.B() {
			return g.controller.PalmLand()
		}
		if oldS.B() && !newS.B() {
			return g.controller.StopLanding()
		}
	} else {
		if !oldS.A() && newS.A() {
			return g.controller.TakeOff()
		}
		if !oldS.B() && newS.B() {
			return g.controller.ThrowTakeOff()
		}
	}

	// handle ox movement
	if oldS.LS().X != newS.LS().X {
		if newS.LS().X >= 0 {
			err = g.controller.Right(newS.LS().X)
		} else {
			err = g.controller.Left(-newS.LS().X)
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

	// handle ox rotation
	if oldS.RS().X != newS.RS().X {
		if newS.RS().X >= 0 {
			err = g.controller.Clockwise(newS.RS().X)
		} else {
			err = g.controller.CounterClockwise(-newS.RS().X)
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

	// handle tricks like flips
	if !oldS.DPadLeft() && newS.DPadLeft() {
		err = g.controller.LeftFlip()
	} else if !oldS.DPadRight() && newS.DPadRight() {
		err = g.controller.RightFlip()
	} else if !oldS.DPadUp() && newS.DPadUp() {
		err = g.controller.FrontFlip()
	} else if !oldS.DPadDown() && newS.DPadDown() {
		err = g.controller.BackFlip()
	} else if !oldS.LB() && newS.LB() {
		err = g.controller.Bounce()
	}
	if err != nil {
		return err
	}

	// handle video settings
	if !oldS.X() && newS.X() {
		e := int(g.info.Exposure()+1) % 3
		err = g.controller.SetExposure(e)
		if err != nil {
			return err
		}
	}
	if !oldS.Y() && newS.Y() {
		r := tello.VideoBitRate(int(g.info.BitRate()+1) % 6)
		err = g.controller.SetVideoBitRate(r)
		if err != nil {
			return err
		}
	}

	// handle speed settings
	if oldS.LT() != 100 && newS.LT() == 100 {
		err = g.controller.SetFastMode()
	}
	if oldS.LT() != -100 && newS.LT() == -100 {
		err = g.controller.SetSlowMode()
	}
	if err != nil {
		return err
	}

	return nil
}