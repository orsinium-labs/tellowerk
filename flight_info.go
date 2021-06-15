package main

import (
	"fmt"

	"gobot.io/x/gobot/platforms/dji/tello"
)

type FlightInfo struct {
	battery  int8
	onGround bool
	exposure int8
}

func NewFlightInfo() FlightInfo {
	return FlightInfo{
		battery:  100,
		onGround: true,
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
		fmt.Println(data.([]byte), int8(data.([]byte)[0]))
		fi.exposure = int8(data.([]byte)[0])
	})
	if err != nil {
		return fmt.Errorf("subscribe to set exposure: %v", err)
	}
	return nil
}

func (fi *FlightInfo) update(data *tello.FlightData) {
	fi.battery = data.BatteryPercentage
	fi.onGround = data.OnGround
}

func (fi FlightInfo) Battery() int8 {
	return fi.battery
}

func (fi FlightInfo) OnGround() bool {
	return fi.onGround
}

func (fi FlightInfo) Exposure() int8 {
	return fi.exposure
}
