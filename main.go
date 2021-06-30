package main

import (
	"fmt"
	"log"

	"github.com/orsinium-labs/gamepad"
	"github.com/orsinium-labs/tellowerk/controllers"
	"github.com/orsinium-labs/tellowerk/plugins"
	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func run(logger *zap.Logger) error {
	var err error

	config := NewConfig()
	driver := tello.NewDriver(fmt.Sprintf("%d", config.Port))

	g, err := gamepad.NewGamepad(config.GamepadID)
	if err != nil {
		return fmt.Errorf("connect gamepad: %v", err)
	}

	// init controller
	controller := controllers.NewMultiplexer()
	controller.Add(controllers.NewLogger(logger))
	if config.Plugins.Driver {
		controller.Add(controllers.NewDriver(driver))
	}

	// init plugins
	pl := plugins.Plugins{
		Controller: controller,
		Logger:     logger,
	}
	if config.Plugins.FFMpeg {
		pl.FFMpeg = plugins.NewFFMpeg(driver)
		pl.FFMpeg.Show = config.Plugins.ImgShow
	}
	if config.Plugins.State {
		pl.State = plugins.NewState(driver)
	}
	if config.Plugins.GamePad {
		pl.GamePad = plugins.NewGamePad(g)
	}
	if config.Plugins.Video {
		pl.Video = plugins.NewVideo(driver)
	}
	if config.Plugins.MPlayer {
		pl.MPlayer = plugins.NewMPlayer(driver)
	}
	if config.Plugins.Recorder {
		pl.Recorder = plugins.NewRecorder(driver)
	}
	if config.Plugins.PiGo {
		pl.PiGo = plugins.NewPiGo()
	}
	if config.Plugins.UI {
		pl.UI = plugins.NewUI()
	}
	pl.State.Addhandler(pl.UI)
	pl.State.Addhandler(plugins.NewStateLogger(logger))

	// start controller
	err = controller.Start()
	if err != nil {
		return fmt.Errorf("start controller: %v", err)
	}
	defer func() {
		err = controller.Stop()
		if err != nil {
			logger.Error("stop controller", zap.Error(err))
		}
	}()

	// run plugins
	err = pl.Run()
	if err != nil {
		return fmt.Errorf("run plugins: %v", err)
	}

	return nil
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("create logger: %v", err)
		return
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	err = run(logger)
	if err != nil {
		logger.Error("runtime error", zap.Error(err))
	}
}
