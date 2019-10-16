package listen

type Ears interface {
	Listen() string
}

func NewEars(engine string, config PocketSphinxConfig) (Ears, error) {
	switch engine {
	case "sphinx", "pocketsphinx":
		ears, err := NewPocketSphinx(config)
		if err != nil {
			ears.Close()
			return nil, err
		}
		return ears, nil
	default:
		panic("unknown engine " + engine)
	}
}
