package plugins

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

const (
	frameX = 480
	frameY = 360
)

type Video struct {
	driver *tello.Driver
	logger *zap.Logger
}

func NewVideo(driver *tello.Driver) *Video {
	return &Video{driver: driver}
}

func (video *Video) Connect(pl *Plugins) {
	video.logger = pl.Logger
}

func (video *Video) Start() error {
	err := video.driver.On(tello.ConnectedEvent, video.configure)
	if err != nil {
		return fmt.Errorf("subscribe to connected: %v", err)
	}
	return nil
}

func (Video) Stop() error {
	return nil
}

func (video *Video) configure(data interface{}) {
	var err error
	video.logger.Debug("connected")
	err = video.driver.StartVideo()
	if err != nil {
		video.logger.Error("start video", zap.Error(err))
		return
	}
	err = video.driver.SetVideoEncoderRate(tello.VideoBitRate1M)
	if err != nil {
		video.logger.Error("set encoder rate", zap.Error(err))
		return
	}
	err = video.driver.SetExposure(0)
	if err != nil {
		video.logger.Error("set exposure", zap.Error(err))
		return
	}
	gobot.Every(1*time.Second, func() {
		err = video.driver.StartVideo()
		if err != nil {
			video.logger.Error("restart video", zap.Error(err))
		}
	})
}
