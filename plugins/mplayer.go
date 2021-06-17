package plugins

import (
	"fmt"
	"io"
	"os/exec"

	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type MPlayer struct {
	driver *tello.Driver
	logger *zap.Logger
	stream io.WriteCloser
}

func NewMPlayer(driver *tello.Driver) *MPlayer {
	return &MPlayer{driver: driver}
}

func (mplayer *MPlayer) Connect(pl *Plugins) {
	mplayer.logger = pl.Logger
}

func (mplayer *MPlayer) Start() error {
	var err error

	// start mplayer subprocess
	cmd := exec.Command("mplayer", "-fps", "30", "-")
	mplayer.stream, err = cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("open stdin pipe to mplayer: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("start mplayer: %v", err)
	}

	// subscribe to events
	err = mplayer.driver.On(tello.VideoFrameEvent, mplayer.handle)
	if err != nil {
		return fmt.Errorf("subscribe to videoframes: %v", err)
	}
	return nil
}

func (mplayer *MPlayer) Stop() error {
	if mplayer.stream != nil {
		err := mplayer.stream.Close()
		if err != nil {
			return fmt.Errorf("close stdin stream: %v", err)
		}
		mplayer.stream = nil
	}
	return nil
}

func (mplayer *MPlayer) handle(data interface{}) {
	if mplayer.stream == nil {
		return
	}
	raw := data.([]byte)
	_, err := mplayer.stream.Write(raw)
	if err != nil {
		mplayer.logger.Error("write frame in stream", zap.Error(err))
	}
}
