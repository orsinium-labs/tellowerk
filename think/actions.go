package think

import (
	"time"

	"github.com/joomcode/errorx"
	"github.com/orsinium/tellowerk/command"
)

func (b *Brain) start(cmd command.Command) (err error) {
	b.logger.Debug("start driver")
	err = b.body.Start()
	if err != nil {
		return errorx.Decorate(err, "cannot start driver")
	}
	b.logger.Debug("start taking off")
	err = b.body.TakeOff()
	if err != nil {
		return errorx.Decorate(err, "cannot take of")
	}
	b.logger.Debug("take off command was sent")
	return nil
}

func (b *Brain) land(cmd command.Command) (err error) {
	b.logger.Debug("start landing")
	err = b.body.Land()
	if err != nil {
		return errorx.Decorate(err, "cannot land")
	}
	b.logger.Debug("landing command was sent")
	return nil
}

func (b *Brain) turnLeft(cmd command.Command) (err error) {
	b.logger.Debug("start rotation")
	var msec time.Duration
	if cmd.Units == command.Seconds {
		msec = time.Duration(cmd.Distance * 1000)
	} else if cmd.Units == command.Degrees {
		msec = time.Duration(cmd.Distance * 12)
	} else {
		msec = 4320 // do 360 degrees rotation by default
	}

	// start rotation
	err = b.body.Clockwise(100)
	if err != nil {
		return errorx.Decorate(err, "cannot start rotation")
	}
	b.logger.Debug("rotation is started")

	// stop rotation
	time.AfterFunc(msec*time.Millisecond, func() {
		b.logger.Debug("stop rotation")
		err = b.body.Clockwise(0)
		if err != nil {
			b.logger.ErrorWith("cannot stop rotation").Err("error", err).Write()
			return
		}
		b.logger.Debug("rotation is stopped")
	})

	return nil
}

func (b *Brain) left(cmd command.Command) (err error) {
	b.logger.DebugWith("start moving").String("direction", "left").Write()
	var msec time.Duration
	if cmd.Units == command.Seconds {
		msec = time.Duration(cmd.Distance * 1000)
	} else if cmd.Units == command.Meters {
		msec = time.Duration(cmd.Distance * 1000)
	} else {
		msec = 1000
	}

	// start moving
	err = b.body.Left(10)
	if err != nil {
		return errorx.Decorate(err, "cannot start moving")
	}
	b.logger.DebugWith("moving started").String("direction", "left").Write()

	// stop moving
	time.AfterFunc(msec*time.Millisecond, func() {
		b.logger.DebugWith("stop moving").String("direction", "left").Write()
		err = b.body.Left(0)
		if err != nil {
			b.logger.ErrorWith("cannot stop rotation").Err("error", err).Write()
			return
		}
		b.logger.DebugWith("moving stopped").String("direction", "left").Write()
	})

	return nil
}
