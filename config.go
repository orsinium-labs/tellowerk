package main

type EnabledPlugins struct {
	// actual plugins
	FFMpeg   bool
	State    bool
	GamePad  bool
	MPlayer  bool
	PiGo     bool
	Video    bool
	Recorder bool
	UI       bool

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
			FFMpeg:   true,
			State:    true,
			GamePad:  true,
			MPlayer:  false,
			PiGo:     false,
			Video:    true,
			Recorder: true,
			UI:       true,

			Driver:    true,
			Targeting: false,
		},
	}
}
