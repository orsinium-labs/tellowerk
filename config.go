package main

type Config struct {
	Port int
	Fly  bool
}

func DefaultConfig() Config {
	return Config{
		Port: 8890,
		Fly:  false,
	}
}
