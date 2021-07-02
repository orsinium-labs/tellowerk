package plugins

import (
	"fmt"

	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type StateHandler interface {
	SetBattery(int8)
	SetWarning(msg string, state bool)
	SetHeight(int16)
	SetNorthSpeed(int16)
	SetEastSpeed(int16)
	SetVerticalSpeed(int16)
}

type State struct {
	logger *zap.Logger
	driver *tello.Driver
	d      tello.FlightData

	exposure int8
	bitrate  tello.VideoBitRate
	face     bool
	photo    bool
	handlers []StateHandler
}

func NewState(driver *tello.Driver) *State {
	return &State{
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
	if fi.d.BatteryPercentage != data.BatteryPercentage {
		for _, h := range fi.handlers {
			h.SetBattery(data.BatteryPercentage)
		}
	}
	if fi.d.EastSpeed != data.EastSpeed {
		for _, h := range fi.handlers {
			h.SetEastSpeed(data.EastSpeed)
		}
	}
	if fi.d.NorthSpeed != data.NorthSpeed {
		for _, h := range fi.handlers {
			h.SetNorthSpeed(data.NorthSpeed)
		}
	}
	if fi.d.VerticalSpeed != data.VerticalSpeed {
		for _, h := range fi.handlers {
			h.SetVerticalSpeed(data.VerticalSpeed)
		}
	}
	if fi.d.Height != data.Height {
		for _, h := range fi.handlers {
			h.SetHeight(data.Height)
		}
	}

	// warnings
	if fi.d.TemperatureHigh != data.TemperatureHigh {
		for _, h := range fi.handlers {
			h.SetWarning("high temperature", data.TemperatureHigh)
		}
	}
	if fi.d.ImuState != data.ImuState {
		for _, h := range fi.handlers {
			h.SetWarning("IMU calibration needed", data.ImuState)
		}
	}
	if fi.d.PressureState != data.PressureState {
		for _, h := range fi.handlers {
			h.SetWarning("pressure issues", data.PressureState)
		}
	}
	if fi.d.PowerState != data.PowerState {
		for _, h := range fi.handlers {
			h.SetWarning("power issues", data.PowerState)
		}
	}
	if fi.d.BatteryState != data.BatteryState {
		for _, h := range fi.handlers {
			h.SetWarning("battery issues", data.BatteryState)
		}
	}
	if fi.d.DownVisualState != data.DownVisualState {
		for _, h := range fi.handlers {
			h.SetWarning("down visibility issues", data.DownVisualState)
		}
	}
	if fi.d.OutageRecording != data.OutageRecording {
		for _, h := range fi.handlers {
			h.SetWarning("video recording outage", data.OutageRecording)
		}
	}
	if fi.d.WindState != data.WindState {
		for _, h := range fi.handlers {
			h.SetWarning("strong wind", data.WindState)
		}
	}

	fi.d = *data
}

func (fi State) Flying() bool {
	return fi.d.Flying
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
