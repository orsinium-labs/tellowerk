package main

import (
	"fmt"

	"gobot.io/x/gobot/platforms/dji/tello"
)

type FlightInfo struct {
	battery  int8
	flying   bool
	exposure int8
	bitrate  tello.VideoBitRate
}

func NewFlightInfo() FlightInfo {
	return FlightInfo{
		battery: 100,
		flying:  false,
	}
}

func (fi *FlightInfo) Subscribe(driver *tello.Driver) error {
	var err error
	err = driver.On(tello.FlightDataEvent, func(data interface{}) {
		fi.update(data.(*tello.FlightData))
	})
	if err != nil {
		return fmt.Errorf("subscribe to flight data: %v", err)
	}
	err = driver.On(tello.SetExposureEvent, func(data interface{}) {
		fi.exposure = int8(data.([]byte)[0])
	})
	if err != nil {
		return fmt.Errorf("subscribe to set exposure: %v", err)
	}
	err = driver.On(tello.SetVideoEncoderRateEvent, func(data interface{}) {
		fi.bitrate = tello.VideoBitRate(data.([]byte)[0])
	})
	if err != nil {
		return fmt.Errorf("subscribe to set bitrate: %v", err)
	}
	return nil
}

func (fi *FlightInfo) update(data *tello.FlightData) {
	fi.battery = data.BatteryPercentage
	fi.flying = data.Flying
}

func (fi FlightInfo) Battery() int8 {
	return fi.battery
}

func (fi FlightInfo) Flying() bool {
	return fi.flying
}

func (fi FlightInfo) Exposure() int8 {
	return fi.exposure
}

func (fi FlightInfo) BitRate() tello.VideoBitRate {
	return fi.bitrate
}
