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

// Config contains settings for Eye
type Config struct {
	Model  string
	FrameX int
	FrameY int
}

// NewEye creates new instance of Eye that can handle video stream from the drone
func NewEye(engine string, config Config, logger *onelog.Logger) (Eye, error) {
	switch engine {
	case "show", "stream", "capture", "mplayer":
		return NewShowEye(logger)
	case "opencv", "gocv", "face":
		return NewOpenCVEye(logger, config)
	default:
		return nil, errors.New("unknown engine: " + engine)
	}
}
