package plugins

import (
	"fmt"

	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type FlightInfo struct {
	logger *zap.Logger
	driver *tello.Driver

	battery  int8
	flying   bool
	exposure int8
	bitrate  tello.VideoBitRate

	// warnings
	temp     bool
	imu      bool
	pressure bool
	video    bool
	wind     bool
}

func NewFlightInfo(driver *tello.Driver) *FlightInfo {
	return &FlightInfo{
		battery: 100,
		driver:  driver,
	}
}

func (fi *FlightInfo) Connect(pl *Plugins) {
	fi.logger = pl.Logger
}

func (FlightInfo) Stop() error {
	return nil
}

func (fi *FlightInfo) Start() error {
	var err error
	err = fi.driver.On(tello.FlightDataEvent, func(data interface{}) {
		fi.update(data.(*tello.FlightData))
	})
	if err != nil {
		return fmt.Errorf("subscribe to flight data: %v", err)
	}
	err = fi.driver.On(tello.SetExposureEvent, func(data interface{}) {
		fi.exposure = int8(data.([]byte)[0])
	})
	if err != nil {
		return fmt.Errorf("subscribe to set exposure: %v", err)
	}
	err = fi.driver.On(tello.SetVideoEncoderRateEvent, func(data interface{}) {
		fi.bitrate = tello.VideoBitRate(data.([]byte)[0])
	})
	if err != nil {
		return fmt.Errorf("subscribe to set bitrate: %v", err)
	}
	return nil
}

func (fi *FlightInfo) update(data *tello.FlightData) {
	if fi.battery != data.BatteryPercentage {
		if data.BatteryPercentage%10 == 0 || fi.battery == 100 {
			field := zap.Int8("value", data.BatteryPercentage)
			if data.BatteryPercentage > 20 {
				fi.logger.Info("battery", field)
			} else {
				fi.logger.Warn("battery", field)
			}
		}
	}
	if !fi.temp && data.TemperatureHigh {
		fi.logger.Warn("high temperature")
	}
	if !fi.imu && data.ImuState {
		fi.logger.Warn("IMU calibration needed")
	}
	if !fi.pressure && data.PressureState {
		fi.logger.Warn("pressure issues")
	}
	if !fi.video && data.OutageRecording {
		fi.logger.Warn("video recording outage")
	}
	if !fi.wind && data.WindState {
		fi.logger.Warn("strong wind")
	}

	fi.temp = data.TemperatureHigh
	fi.imu = data.ImuState
	fi.pressure = data.PressureState
	fi.video = data.OutageRecording
	fi.wind = data.WindState

	fi.battery = data.BatteryPercentage
	fi.flying = data.Flying
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
