package speak

import "errors"

type Voice interface {
	Say(text string) error
}

func NewVoice(engine string, speaker string) (Voice, error) {
	switch engine {
	case "flite":
		return NewFlite(speaker)
	default:
		return nil, errors.New("unknown engine: " + engine)
	}
}
