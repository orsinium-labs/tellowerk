package see

import (
	"io"
	"os/exec"

	"github.com/joomcode/errorx"

	"github.com/francoispqt/onelog"
)

// ShowEye translates video from drone on the laptop screen
type ShowEye struct {
	stream io.WriteCloser
	logger *onelog.Logger
}

// Handle processes a new frame on a screen
func (eye *ShowEye) Handle(data interface{}) {
	if eye.stream == nil {
		return
	}
	raw := data.([]byte)
	_, err := eye.stream.Write(raw)
	if err != nil {
		eye.logger.ErrorWith("cannot write frame in stream").Err("error", err).Write()
	}
}

// Close closes stdin stream to mlayer
func (eye *ShowEye) Close() (err error) {
	if eye.stream != nil {
		err = eye.stream.Close()
		eye.stream = nil
	}
	return err
}

// NewShowEye creates ShowEye instance to stream video from drone on a screen
func NewShowEye(logger *onelog.Logger) (*ShowEye, error) {
	var err error
	eye := ShowEye{logger: logger}

	mplayer := exec.Command("mplayer", "-fps", "20", "-")
	eye.stream, err = mplayer.StdinPipe()
	if err != nil {
		return &eye, errorx.Decorate(err, "cannot open stdin pipe to mplayer")
	}

	err = mplayer.Start()
	if err != nil {
		return &eye, errorx.Decorate(err, "cannot start mplayer")
	}
	return &eye, nil
}
