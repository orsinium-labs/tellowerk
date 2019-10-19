package think

import (
	"errors"

	"github.com/francoispqt/onelog"
	"github.com/orsinium/tellowerk/command"
	"gobot.io/x/gobot/platforms/dji/tello"
)

// Brain accepts command and says to body what to do.
type Brain struct {
	body   *tello.Driver
	logger *onelog.Logger
	dry    bool
}

// Do does actions for given command
func (b *Brain) Do(cmd command.Command) error {
	b.logger.DebugWith("start doing action").String("action", string(cmd.Action)).Write()
	if b.dry {
		b.logger.Debug("dry run")
		return nil
	}
	switch cmd.Action {
	case command.Start:
		return b.start(cmd)
	case command.Land:
		return b.land(cmd)
	case command.TurnLeft:
		return b.turnLeft(cmd)
	}
	return errors.New("unknown command")
}

// Stop stops the driver
func (b *Brain) Stop() error {
	b.logger.Debug("stopping the driver")
	if b.dry {
		b.logger.Debug("dry run")
		return nil
	}
	err := b.body.Halt()
	if err != nil {
		return err
	}
	b.logger.Debug("command to stop the driver was sent")
	return nil
}

// NewBrain creates Brain instance to do actions
func NewBrain(dry bool, body *tello.Driver, logger *onelog.Logger) *Brain {
	return &Brain{dry: dry, body: body, logger: logger}
}
