package listen

import (
	"errors"

	"github.com/francoispqt/onelog"
)

type Ear interface {
	Listen() string
	Close() error
}

type ListenConfig struct {
	HMM        string
	Dict       string
	LM         string
	Samples    int
	SampleRate int `toml:"sample_rate"`
	Channels   int
}

func NewEar(engine string, config ListenConfig, logger *onelog.Logger) (Ear, error) {
	switch engine {
	case "sphinx", "pocketsphinx":
		return NewPocketSphinx(config, logger)
	default:
		return nil, errors.New("unknown engine: " + engine)
	}
}
