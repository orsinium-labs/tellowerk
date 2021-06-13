package main

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func run(logger *zap.Logger) error {
	var err error

	config := DefaultConfig()
	driver := tello.NewDriver(fmt.Sprintf("%d", config.Port))

	finfo := FlightInfo{}
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

	controller := Controller{
		driver: driver,
		logger: logger,
		fly:    config.Fly,
	}
	err = controller.Start()
	if err != nil {
		return fmt.Errorf("start controller: %v", err)
	}
	defer func() {
		err = controller.Stop()
		if err != nil {
			logger.Error("cannot stop controller", zap.Error(err))
		}
	}()

	err = controller.TakeOff()
	if err != nil {
		return fmt.Errorf("take off: %v", err)
	}
	time.Sleep(4 * time.Second)
	err = controller.Clockwise()
	if err != nil {
		return fmt.Errorf("clockwise: %v", err)
	}
	time.Sleep(10 * time.Second)
	err = controller.Land()
	if err != nil {
		return fmt.Errorf("land: %v", err)
	}
	logger.Info("battery", zap.Int8("value", finfo.Battery()))

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
