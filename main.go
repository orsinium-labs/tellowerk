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

	controller := Controller{
		driver: tello.NewDriver(fmt.Sprintf("%d", config.Port)),
		logger: logger,
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

	finfo := FlightInfo{}
	err = finfo.Subscribe(controller.driver)
	if err != nil {
		return fmt.Errorf("subscribe to flight info: %v", err)
	}

	err = controller.TakeOff()
	if err != nil {
		return fmt.Errorf("take off: %v", err)
	}
	time.Sleep(4 * time.Second)
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
