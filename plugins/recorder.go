package plugins

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Recorder struct {
	driver *tello.Driver
	logger *zap.Logger
	stream io.WriteCloser
}

func NewRecorder(driver *tello.Driver) *Recorder {
	return &Recorder{driver: driver}
}

func (rec *Recorder) Connect(pl *Plugins) {
	rec.logger = pl.Logger
}

func (rec *Recorder) Start() error {
	var err error
	fname := fmt.Sprintf("tello-%s.mpg", time.Now().Format("2006-01-02_15-04-05"))
	rec.stream, err = os.Create(fname)
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}
	err = rec.driver.On(tello.VideoFrameEvent, rec.handle)
	if err != nil {
		return fmt.Errorf("subscribe to video frames: %v", err)
	}
	return nil
}

func (rec *Recorder) Stop() error {
	if rec.stream != nil {
		err := rec.stream.Close()
		if err != nil {
			return fmt.Errorf("close stdin stream: %v", err)
		}
		rec.stream = nil
	}
	return nil
}

func (rec *Recorder) handle(data interface{}) {
	if rec.stream == nil {
		return
	}
	raw := data.([]byte)
	_, err := rec.stream.Write(raw)
	if err != nil {
		rec.logger.Error("write frame in stream", zap.Error(err))
	}
}
