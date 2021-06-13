package act

import (
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Body interface {
	TakeOff() error
	Land() error
	Start() error
}

func NewBody() *tello.Driver {
	return tello.NewDriver("8890")
	// drone.Start()
}
