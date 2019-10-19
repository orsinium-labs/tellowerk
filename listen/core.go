package listen

import (
	"errors"

	"github.com/francoispqt/onelog"
)

// Ear is a common interface for all ears
type Ear interface {
	Listen() string
	Close() error
}

// ListenConfig is a format of settings for ears in config.toml
type ListenConfig struct {
	HMM        string
	Dict       string
	LM         string
	Samples    int
	SampleRate int `toml:"sample_rate"`
	Channels   int
	Scenario   string
}

// NewEar creates new instance of Ear that can listen for user commands in some way
func NewEar(engine string, config ListenConfig, logger *onelog.Logger) (Ear, error) {
	switch engine {
	case "sphinx", "pocketsphinx":
		return NewPocketSphinx(config, logger)
	case "stdin", "console", "terminal":
		return NewConsole(config, logger)
	case "scenario", "file":
		return NewScenario(config, logger)
	default:
		return nil, errors.New("unknown engine: " + engine)
	}
}
