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

type Video struct {
	driver *tello.Driver
	logger *zap.Logger
	stream io.WriteCloser
}

func (v *Video) Start() error {
	var err error

	// start mplayer subprocess
	mplayer := exec.Command("mplayer", "-fps", "25", "-")
	v.stream, err = mplayer.StdinPipe()
	if err != nil {
		return fmt.Errorf("open stdin pipe to mplayer: %v", err)
	}
	err = mplayer.Start()
	if err != nil {
		return fmt.Errorf("start mplayer: %v", err)
	}

	// subscribe to events
	err = v.driver.On(tello.ConnectedEvent, v.configure)
	if err != nil {
		return fmt.Errorf("subscribe to connected: %v", err)
	}
	err = v.driver.On(tello.VideoFrameEvent, v.handle)
	if err != nil {
		return fmt.Errorf("subscribe to videoframes: %v", err)
	}
	return nil
}

func (v *Video) Stop() error {
	if v.stream != nil {
		err := v.stream.Close()
		if err != nil {
			return fmt.Errorf("close stdin stream: %v", err)
		}
		v.stream = nil
	}
	return nil
}

func (v *Video) configure(data interface{}) {
	var err error
	v.logger.Debug("connected")
	err = v.driver.StartVideo()
	if err != nil {
		v.logger.Error("start video", zap.Error(err))
		return
	}
	err = v.driver.SetVideoEncoderRate(tello.VideoBitRate1M)
	if err != nil {
		v.logger.Error("set encoder rate", zap.Error(err))
		return
	}
	err = v.driver.SetExposure(0)
	if err != nil {
		v.logger.Error("set exposure", zap.Error(err))
		return
	}
	gobot.Every(1*time.Second, func() {
		err = v.driver.StartVideo()
		if err != nil {
			v.logger.Error("restart video", zap.Error(err))
		}
	})
}

func (eye *Video) handle(data interface{}) {
	if eye.stream == nil {
		return
	}
	raw := data.([]byte)
	_, err := eye.stream.Write(raw)
	if err != nil {
		eye.logger.Error("write frame in stream", zap.Error(err))
	}
}
