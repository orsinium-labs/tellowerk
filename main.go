package main

import (
	"fmt"
	"log"
	"time"

	"github.com/orsinium-labs/gamepad"
	"github.com/orsinium-labs/tellowerk/controllers"
	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func run(logger *zap.Logger) error {
	var err error

	config := DefaultConfig()
	driver := tello.NewDriver(fmt.Sprintf("%d", config.Port))

	finfo := NewFlightInfo()
	err = finfo.Subscribe(driver)
	if err != nil {
		return fmt.Errorf("subscribe to flight info: %v", err)
	}

	mplayer := MPlayer{
		driver: driver,
		logger: logger,
	}
	err = mplayer.Start()
	if err != nil {
		return fmt.Errorf("start video: %v", err)
	}
	defer func() {
		err = mplayer.Stop()
		if err != nil {
			logger.Error("cannot stop video", zap.Error(err))
		}
	}()

	controller := controllers.NewMultiplexer()
	controller.Add(controllers.NewLogger(logger))
	if config.Fly {
		controller.Add(controllers.NewDriver(driver))
	}
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

	finish := make(chan struct{})
	g, err := gamepad.NewGamepad(config.GamepadID)
	if err != nil {
		return fmt.Errorf("connect gamepad: %v", err)
	}
	gamepad := GamePad{
		controller: controller,
		info:       &finfo,
		gamepad:    g,
		logger:     logger,
		finish:     finish,
	}
	gamepad.Start()

	time.Sleep(time.Second)
	logger.Info("battery", zap.Int8("value", finfo.Battery()))
	<-finish

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
