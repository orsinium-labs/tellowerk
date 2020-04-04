package think

import (
	"context"
	"time"

	"github.com/joomcode/errorx"
	"github.com/orsinium-labs/tellowerk/command"
)

// after is a blocking wrapper that calls all `cancels`
// and then  calls the `action` when given `duration` expires
// or if `cancel` from `cancels` is called.
func (b *Brain) after(duration time.Duration, action func()) {
	b.cancel()
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	b.logger.Debug("pushing action")
	b.cancels <- cancel
	b.logger.Debug("action has been pushed")
	<-ctx.Done()
	b.logger.Debug("an action is called")
	action()
}

func (b *Brain) cancel() {
	for {
		select {
		case cancel := <-b.cancels:
			b.logger.Debug("call an old action")
			cancel()
		default:
			return
		}
	}
}

// up and down //

func (b *Brain) start(cmd command.Command) (err error) {
	b.logger.DebugWith("start taking off").Int("job", cmd.JobID).Write()
	err = b.body.TakeOff()
	if err != nil {
		return errorx.Decorate(err, "cannot take of")
	}
	b.logger.DebugWith("take off command was sent").Int("job", cmd.JobID).Write()
	return nil
}

func (b *Brain) land(cmd command.Command) (err error) {
	b.cancel()
	b.logger.DebugWith("start landing").Int("job", cmd.JobID).Write()
	err = b.body.Land()
	if err != nil {
		return errorx.Decorate(err, "cannot land")
	}
	b.logger.DebugWith("landing command was sent").Int("job", cmd.JobID).Write()
	return nil
}

// rotation oxy //

func (b *Brain) turnLeft(cmd command.Command) (err error) {
	return b.turn(b.body.CounterClockwise, "left", cmd)
}

func (b *Brain) turnRight(cmd command.Command) (err error) {
	return b.turn(b.body.Clockwise, "right", cmd)
}

func (b *Brain) turn(handler func(int) error, direction string, cmd command.Command) (err error) {
	b.logger.DebugWith("start rotation").String("direction", direction).Int("job", cmd.JobID).Write()
	var msec time.Duration
	if cmd.Units == command.Seconds {
		msec = time.Duration(cmd.Distance * 1000)
	} else if cmd.Units == command.Degrees {
		msec = time.Duration(cmd.Distance * 11)
	} else {
		msec = 1000 // do 90 degrees rotation by default
	}

	// start rotation
	err = handler(100)
	if err != nil {
		return errorx.Decorate(err, "cannot start rotation")
	}
	b.logger.DebugWith("rotation is started").String("direction", direction).Int("job", cmd.JobID).Write()

	// stop rotation
	go b.after(msec*time.Millisecond, func() {
		b.logger.DebugWith("stop rotation").String("direction", direction).Int("job", cmd.JobID).Write()
		err = b.body.Clockwise(0)
		if err != nil {
			b.logger.ErrorWith("cannot stop rotation").Err("error", err).Int("job", cmd.JobID).Write()
			return
		}
		b.logger.DebugWith("rotation is stopped").String("direction", direction).Int("job", cmd.JobID).Write()
	})

	return nil
}

// movement oxy //

func (b *Brain) left(cmd command.Command) error {
	return b.move(b.body.Left, "left", cmd)
}

func (b *Brain) right(cmd command.Command) error {
	return b.move(b.body.Right, "right", cmd)
}

func (b *Brain) forward(cmd command.Command) error {
	return b.move(b.body.Forward, "forward", cmd)
}

func (b *Brain) backward(cmd command.Command) error {
	return b.move(b.body.Backward, "backward", cmd)
}

func (b *Brain) up(cmd command.Command) error {
	return b.move(b.body.Up, "up", cmd)
}

func (b *Brain) down(cmd command.Command) error {
	return b.move(b.body.Down, "down", cmd)
}

func (b *Brain) move(handler func(int) error, direction string, cmd command.Command) (err error) {
	b.logger.DebugWith("start moving").String("direction", direction).Int("job", cmd.JobID).Write()
	var msec time.Duration
	if cmd.Units == command.Seconds {
		msec = time.Duration(cmd.Distance * 1000)
	} else if cmd.Units == command.Meters {
		msec = time.Duration(cmd.Distance*1000 - 50)
	} else {
		msec = 950
	}

	// start moving
	err = handler(50)
	if err != nil {
		return errorx.Decorate(err, "cannot start moving")
	}
	b.logger.DebugWith("moving started").String("direction", direction).Int("job", cmd.JobID).Write()

	// stop moving
	go b.after(msec*time.Millisecond, func() {
		b.logger.DebugWith("stop moving").String("direction", direction).Int("job", cmd.JobID).Write()
		err = handler(0)
		if err != nil {
			b.logger.ErrorWith("cannot stop rotation").Err("error", err).Int("job", cmd.JobID).Write()
			return
		}
		b.logger.DebugWith("moving stopped").String("direction", direction).Int("job", cmd.JobID).Write()
	})

	return nil
}
