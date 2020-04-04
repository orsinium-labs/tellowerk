package think

import (
	"errors"

	"github.com/francoispqt/onelog"
	"github.com/orsinium-labs/tellowerk/command"
	"gobot.io/x/gobot/platforms/dji/tello"
)

// Brain accepts command and says to body what to do.
type Brain struct {
	body     *tello.Driver
	logger   *onelog.Logger
	dry      bool
	registry map[command.Action]func(cmd command.Command) error
	delayed  DelayedCalls
}

// Do does actions for given command
func (b *Brain) Do(cmd command.Command) error {
	b.logger.DebugWith("start doing action").String("action", string(cmd.Action)).Int("job", cmd.JobID).Write()
	if b.dry {
		b.logger.Debug("dry run")
		return nil
	}
	handler, ok := b.registry[cmd.Action]
	if !ok {
		return errors.New("unknown command")
	}
	return handler(cmd)
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

func (b *Brain) register(cmd command.Action, handler func(cmd command.Command) error) {
	b.registry[cmd] = handler
}

// NewBrain creates Brain instance to do actions
func NewBrain(dry bool, body *tello.Driver, logger *onelog.Logger) *Brain {
	b := Brain{
		dry:      dry,
		body:     body,
		logger:   logger,
		registry: make(map[command.Action]func(cmd command.Command) error),
		delayed: DelayedCalls{
			queue:  make(chan *DelayedCall, 100),
			logger: logger,
		},
	}

	b.register(command.Start, b.start)
	b.register(command.Land, b.land)

	b.register(command.Left, b.left)
	b.register(command.Right, b.right)
	b.register(command.Front, b.forward)
	b.register(command.Back, b.backward)

	b.register(command.Up, b.up)
	b.register(command.Down, b.down)

	b.register(command.TurnLeft, b.turnLeft)
	b.register(command.TurnRight, b.turnRight)

	return &b
}
