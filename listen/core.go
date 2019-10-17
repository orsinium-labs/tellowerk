package listen

import (
	"errors"
)

type Ear interface {
	Listen() string
}

type ListenConfig struct {
	HMM        string
	Dict       string
	LM         string
	Samples    int
	SampleRate int `toml:"sample_rate"`
	Channels   int
}

func NewEar(engine string, config ListenConfig) (Ear, error) {
	switch engine {
	case "sphinx", "pocketsphinx":
		ears, err := NewPocketSphinx(config)
		if err != nil {
			ears.Close()
			return nil, err
		}
		return ears, nil
	default:
		return nil, errors.New("unknown engine: " + engine)
	}
}
