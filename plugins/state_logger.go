package plugins

import "go.uber.org/zap"

type StateLogger struct {
	logger  *zap.Logger
	battery bool
}

var _ StateHandler = &StateLogger{}

func (log *StateLogger) SetBattery(val int8) {
	if val%10 == 0 || !log.battery {
		log.battery = true
		field := zap.Int8("value", val)
		if val > 20 {
			log.logger.Info("battery", field)
		} else {
			log.logger.Warn("battery", field)
		}
	}
}

func (log *StateLogger) SetWarning(msg string, state bool) {
	if state {
		log.logger.Warn(msg)
	}
}
