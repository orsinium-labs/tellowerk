module github.com/orsinium-labs/tellowerk

go 1.16

replace (
	github.com/orsinium-labs/gamepad => ../gamepad
)

require (
	github.com/orsinium-labs/gamepad v1.0.4
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0
	gobot.io/x/gobot v1.15.0
)
