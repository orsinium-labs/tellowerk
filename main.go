package main

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/francoispqt/onelog"
	"github.com/orsinium/tellowerk/listen"
	"github.com/orsinium/tellowerk/speak"
)

type configListen struct {
	listen.ListenConfig
	Engine string
}
type configSpeak struct {
	Engine  string
	Speaker string
}

type Config struct {
	Listen configListen
	Speak  configSpeak
}

func main() {
	logger := onelog.New(os.Stdout, onelog.ALL)

	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		logger.FatalWith("cannot read config").Err("error", err).Write()
	}

	ears, err := listen.NewEar(conf.Listen.Engine, conf.Listen.ListenConfig)
	if err != nil {
		logger.FatalWith("cannot make ear").Err("error", err).Write()
	}

	voice, err := speak.NewVoice(conf.Speak.Engine, conf.Speak.Speaker)
	if err != nil {
		logger.FatalWith("cannot make voice").Err("error", err).Write()
	}

	logger.Info("start")
	for {
		text := ears.Listen()
		logger.InfoWith("text heared").String("text", text).Write()
		err = voice.Say(text)
		if err != nil {
			logger.ErrorWith("cannot say").Err("error", err).Write()
		}
	}
}
