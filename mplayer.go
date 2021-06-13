package main

import (
	"fmt"
	"io"
	"os/exec"
	"time"

	"go.uber.org/zap"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type MPlayer struct {
	driver *tello.Driver
	logger *zap.Logger
	stream io.WriteCloser
}

func (mplayer *MPlayer) Start() error {
	var err error

	// start mplayer subprocess
	cmd := exec.Command("mplayer", "-fps", "25", "-")
	mplayer.stream, err = cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("open stdin pipe to mplayer: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("start mplayer: %v", err)
	}

	// subscribe to events
	err = mplayer.driver.On(tello.ConnectedEvent, mplayer.configure)
	if err != nil {
		return fmt.Errorf("subscribe to connected: %v", err)
	}
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

func (mplayer *MPlayer) configure(data interface{}) {
	var err error
	mplayer.logger.Debug("connected")
	err = mplayer.driver.StartVideo()
	if err != nil {
		mplayer.logger.Error("start video", zap.Error(err))
		return
	}
	err = mplayer.driver.SetVideoEncoderRate(tello.VideoBitRate1M)
	if err != nil {
		mplayer.logger.Error("set encoder rate", zap.Error(err))
		return
	}
	err = mplayer.driver.SetExposure(0)
	if err != nil {
		mplayer.logger.Error("set exposure", zap.Error(err))
		return
	}
	gobot.Every(1*time.Second, func() {
		err = mplayer.driver.StartVideo()
		if err != nil {
			mplayer.logger.Error("restart video", zap.Error(err))
		}
	})
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
