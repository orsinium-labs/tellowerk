package plugins

import (
	"time"

	"github.com/orsinium-labs/gamepad"
	"github.com/orsinium-labs/tellowerk/controllers"
	"go.uber.org/zap"
)

type GamePad struct {
	controller controllers.Controller
	state      *State
	ffmpeg     *FFMpeg
	logger     *zap.Logger

	gamepad *gamepad.GamePad
	ui      *UI
}

func NewGamePad(g *gamepad.GamePad) *GamePad {
	return &GamePad{
		gamepad: g,
	}
}

func (g *GamePad) Connect(pl *Plugins) {
	g.controller = pl.Controller
	g.state = pl.State
	g.logger = pl.Logger
	g.ffmpeg = pl.FFMpeg
}

func (g *GamePad) Start() error {
	go g.worker()
	return nil
}

func (GamePad) Stop() error {
	return nil
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
			err = g.ui.Stop()
			if err != nil {
				g.logger.Error("stop UI", zap.Error(err))
			}
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

	// take off and land
	if g.state.Flying() {
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

	// movement
	if oldS.LS().X != newS.LS().X {
		err = g.controller.OX(newS.LS().X)
		if err != nil {
			return err
		}
	}
	if oldS.LS().Y != newS.LS().Y {
		err = g.controller.OY(newS.LS().Y)
		if err != nil {
			return err
		}
	}
	if oldS.RS().X != newS.RS().X {
		err = g.controller.Rotate(newS.RS().X)
		if err != nil {
			return err
		}
	}
	if oldS.RS().Y != newS.RS().Y {
		err = g.controller.OZ(-newS.RS().Y)
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
		e := int(g.state.Exposure()+1) % 3
		err = g.controller.SetExposure(e)
		if err != nil {
			return err
		}
	}
	if !oldS.Y() && newS.Y() {
		// r := tello.VideoBitRate(int(g.state.BitRate()+1) % 6)
		// err = g.controller.SetVideoBitRate(r)
		g.state.SetFaceCapture(!g.state.FaceCapture())
	}
	// take a photo
	if !oldS.Guide() && newS.Guide() && g.ffmpeg != nil {
		g.state.SetTakePhoto(true)
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
