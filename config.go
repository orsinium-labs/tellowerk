package main

type EnabledPlugins struct {
	// actual plugins
	FFMpeg     bool
	FlightInfo bool
	GamePad    bool
	MPlayer    bool
	PiGo       bool
	Video      bool

	// subplugins
	Driver    bool
	Targeting bool
}

type Config struct {
	Port      int
	GamepadID int
	Plugins   EnabledPlugins
}

func NewConfig() Config {
	return Config{
		Port:      8890,
		GamepadID: 1,
		Plugins: EnabledPlugins{
			FFMpeg:     false,
			FlightInfo: true,
			GamePad:    true,
			MPlayer:    true,
			PiGo:       false,
			Video:      true,

			Driver:    true,
			Targeting: true,
		},
	}
}
