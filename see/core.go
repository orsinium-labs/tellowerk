package see

import (
	"errors"

	"github.com/francoispqt/onelog"
)

// Eye is an interface for video stream handlers
type Eye interface {
	Handle(interface{})
	Close() error
}

// NewEye creates new instance of Eye that can handle video stream from the drone
func NewEye(engine string, logger *onelog.Logger) (Eye, error) {
	switch engine {
	case "show", "stream", "capture", "mplayer":
		return NewShowEye(logger)
	default:
		return nil, errors.New("unknown engine: " + engine)
	}
}
