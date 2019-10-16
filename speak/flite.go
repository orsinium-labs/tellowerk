package speak

import "github.com/gen2brain/flite-go"

type Flite struct {
	voice *flite.Voice
}

func (f *Flite) Say(text string) error {
	flite.TextToSpeech(text, f.voice, "play")
	return nil
}

func NewFlite(speaker string) *Flite {
	voice, err := flite.VoiceSelect(speaker)
	if err != nil {
		panic(err)
	}
	return &Flite{voice: voice}
}
