package think

import (
	"time"

	"github.com/francoispqt/onelog"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func NewBrain(body *tello.Driver, logger *onelog.Logger) *gobot.Robot {
	work := func() {
		logger.Info("start work")
		err := body.TakeOff()
		if err != nil {
			logger.ErrorWith("cannot take off").Err("error", err).Write()
		}

		gobot.After(2*time.Second, func() {
			logger.Info("rotate")
			body.Clockwise(100)
			logger.Info("end rotate")
		})

		gobot.After(6*time.Second, func() {
			logger.Info("land")
			body.Clockwise(0)
			err := body.Land()
			if err != nil {
				logger.ErrorWith("cannot land").Err("error", err).Write()
			}
			logger.Info("end land")
		})
		logger.Info("end work")
	}

	// body.Start()

	return gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{body},
		work,
	)
}
