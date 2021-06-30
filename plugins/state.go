package plugins

import (
	"fmt"

	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type StateHandler interface {
	SetBattery(int8)
	SetWarning(msg string, state bool)

	SetNorthSpeed(int16)
	SetEastSpeed(int16)
	SetVerticalSpeed(int16)
}

type State struct {
	logger *zap.Logger
	driver *tello.Driver

	flying   bool
	exposure int8
	bitrate  tello.VideoBitRate
	face     bool
	photo    bool
	handlers []StateHandler

	// warnings
	temp     bool
	imu      bool
	pressure bool
	video    bool
	wind     bool

	// metrics
	battery int8
	east    int16
	north   int16
	vert    int16
}

func NewState(driver *tello.Driver) *State {
	return &State{
		battery:  100,
		driver:   driver,
		handlers: make([]StateHandler, 0),
	}
}

func (fi *State) Connect(pl *Plugins) {
	fi.logger = pl.Logger
}

func (fi *State) Addhandler(h StateHandler) {
	fi.handlers = append(fi.handlers, h)
}

func (State) Stop() error {
	return nil
}

func (fi *State) Start() error {
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

func (fi *State) update(data *tello.FlightData) {
	// metrics
	if fi.battery != data.BatteryPercentage {
		for _, h := range fi.handlers {
			h.SetBattery(data.BatteryPercentage)
		}
		fi.battery = data.BatteryPercentage
	}
	if fi.east != data.EastSpeed {
		for _, h := range fi.handlers {
			h.SetEastSpeed(data.EastSpeed)
		}
		fi.east = data.EastSpeed
	}
	if fi.north != data.NorthSpeed {
		for _, h := range fi.handlers {
			h.SetNorthSpeed(data.NorthSpeed)
		}
		fi.north = data.NorthSpeed
	}
	if fi.vert != data.VerticalSpeed {
		for _, h := range fi.handlers {
			h.SetVerticalSpeed(data.VerticalSpeed)
		}
		fi.vert = data.VerticalSpeed
	}

	// warnings
	if fi.temp != data.TemperatureHigh {
		for _, h := range fi.handlers {
			h.SetWarning("high temperature", data.TemperatureHigh)
		}
		fi.temp = data.TemperatureHigh
	}
	if fi.imu != data.ImuState {
		for _, h := range fi.handlers {
			h.SetWarning("IMU calibration needed", data.ImuState)
		}
		fi.imu = data.ImuState
	}
	if fi.pressure != data.PressureState {
		for _, h := range fi.handlers {
			h.SetWarning("pressure issues", data.PressureState)
		}
		fi.pressure = data.PressureState
	}
	if fi.video != data.OutageRecording {
		for _, h := range fi.handlers {
			h.SetWarning("video recording outage", data.OutageRecording)
		}
		fi.video = data.OutageRecording
	}
	if fi.wind != data.WindState {
		for _, h := range fi.handlers {
			h.SetWarning("strong wind", data.WindState)
		}
		fi.wind = data.WindState
	}

	fi.flying = data.Flying
}

func (fi State) Flying() bool {
	return fi.flying
}

func (fi State) Exposure() int8 {
	return fi.exposure
}

func (fi State) BitRate() tello.VideoBitRate {
	return fi.bitrate
}

func (fi State) FaceCapture() bool {
	return fi.face
}

func (fi *State) SetFaceCapture(val bool) {
	fi.logger.Debug("set face capture", zap.Bool("val", val))
	fi.photo = val
}

func (fi State) TakePhoto() bool {
	return fi.face
}

func (fi *State) SetTakePhoto(val bool) {
	fi.logger.Debug("set take photo", zap.Bool("val", val))
	fi.photo = val
}
