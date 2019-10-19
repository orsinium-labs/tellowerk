package main

import (
	"os"
	"time"

	"github.com/orsinium/tellowerk/command"

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

// Config is a storage for all app settings read from config.toml
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

	ear, err := listen.NewEar(conf.Listen.Engine, conf.Listen.ListenConfig, logger)
	defer func() {
		ear.Close()
		if err != nil {
			logger.ErrorWith("cannot stop listening").Err("error", err).Write()
		}
	}()
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
	defer func() {
		err = brain.Stop()
		if err != nil {
			logger.ErrorWith("cannot stop driver").Err("error", err).Write()
		}
	}()

	logger.Info("start")
	for {
		text := ear.Listen()
		logger.DebugWith("text heared").String("text", text).Write()
		cmd := command.Understand(text)
		if cmd.Action == "" {
			logger.DebugWith("cannot recognize command").String("text", text).Write()
			continue
		}
		err = brain.Do(cmd)
		if err != nil {
			logger.ErrorWith("cannot do action").Err("error", err).Write()
			continue
		}
		err = voice.Say(string(cmd.Action))
		if err != nil {
			logger.ErrorWith("cannot say").Err("error", err).Write()
		}
	}
}
