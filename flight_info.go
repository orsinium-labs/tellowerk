package main

import "gobot.io/x/gobot/platforms/dji/tello"

type FlightInfo struct {
	battery  int8
	onGround bool
}

func NewFlightInfo() FlightInfo {
	return FlightInfo{
		battery:  100,
		onGround: true,
	}
}

func (fi *FlightInfo) Subscribe(driver *tello.Driver) error {
	return driver.On(tello.FlightDataEvent, func(data interface{}) {
		fi.update(data.(*tello.FlightData))
	})
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
