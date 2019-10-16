package speak

type Voice interface {
	Say(text string) error
}

func NewVoice(engine string, speaker string) Voice {
	switch engine {
	case "flite":
		return NewFlite(speaker)
	default:
		panic("unknown engine " + engine)
	}
}
