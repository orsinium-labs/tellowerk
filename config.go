package main

type Config struct {
	Port      int
	GamepadID int
	Fly       bool
}

func DefaultConfig() Config {
	return Config{
		Port:      8890,
		GamepadID: 1,
		Fly:       false,
	}
}
