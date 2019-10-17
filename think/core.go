package think

import (
	"time"

	"github.com/francoispqt/onelog"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func NewBrain(body *tello.Driver, logger *onelog.Logger) *gobot.Robot {
	work := func() {
		err := body.TakeOff()
		if err != nil {
			logger.ErrorWith("cannot take off").Err("error", err).Write()
		}

		gobot.After(5*time.Second, func() {
			err := body.Land()
			if err != nil {
				logger.ErrorWith("cannot land").Err("error", err).Write()
			}
		})
	}

	// body.Start()

	return gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{body},
		work,
	)
}
