package speak

import "github.com/gen2brain/flite-go"

type Flite struct {
	voice *flite.Voice
}

func (f *Flite) Say(text string) error {
	flite.TextToSpeech(text, f.voice, "play")
	return nil
}

func NewFlite(speaker string) (*Flite, error) {
	voice, err := flite.VoiceSelect(speaker)
	if err != nil {
		return nil, err
	}
	return &Flite{voice: voice}, nil
}
