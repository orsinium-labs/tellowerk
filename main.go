package main

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/francoispqt/onelog"
	"github.com/orsinium/tellowerk/act"
	"github.com/orsinium/tellowerk/listen"
	"github.com/orsinium/tellowerk/speak"
	"github.com/orsinium/tellowerk/think"
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
	logger.Hook(func(e onelog.Entry) {
		e.String("time", time.Now().Format("15:04:05"))
	})

	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		logger.FatalWith("cannot read config").Err("error", err).Write()
	}

	ears, err := listen.NewEar(conf.Listen.Engine, conf.Listen.ListenConfig, logger)
	if err != nil {
		logger.FatalWith("cannot make ear").Err("error", err).Write()
	}

	voice, err := speak.NewVoice(conf.Speak.Engine, conf.Speak.Speaker)
	if err != nil {
		logger.FatalWith("cannot make voice").Err("error", err).Write()
	}

	logger.Info("start thinking")
	body := act.NewBody()
	brain := think.NewBrain(body, logger)
	brain.Start()
	defer brain.Stop()
	// defer body.Halt()

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
