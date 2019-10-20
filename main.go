package main

import (
	"os"
	"time"

	"github.com/joomcode/errorx"
	"github.com/orsinium/tellowerk/command"
	"github.com/orsinium/tellowerk/see"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"

	"github.com/BurntSushi/toml"
	"github.com/francoispqt/onelog"
	"github.com/orsinium/tellowerk/act"
	"github.com/orsinium/tellowerk/listen"
	"github.com/orsinium/tellowerk/speak"
	"github.com/orsinium/tellowerk/think"
)

type configListen struct {
	listen.Config
	Engine string
}
type configThink struct {
	Dry bool
}
type configSpeak struct {
	Engine  string
	Speaker string
}

type configSee struct {
	see.Config
	Engine string
}

// Config is a storage for all app settings read from config.toml
type Config struct {
	Listen configListen
	Speak  configSpeak
	Think  configThink
	See    configSee
}

func start(logger *onelog.Logger, body *tello.Driver, eye see.Eye) (err error) {
	logger.Debug("registering connection handler")
	err = body.On(tello.ConnectedEvent, func(data interface{}) {
		logger.Debug("connected")
		body.StartVideo()
		body.SetVideoEncoderRate(tello.VideoBitRateAuto)
		body.SetExposure(0)
		gobot.Every(100*time.Millisecond, func() {
			body.StartVideo()
		})
	})

	logger.Debug("registering video handler")
	err = body.On(tello.VideoFrameEvent, eye.Handle)
	if err != nil {
		return errorx.Decorate(err, "cannot register handler for video")
	}

	logger.Debug("starting driver")
	err = body.Start()
	if err != nil {
		return errorx.Decorate(err, "cannot start driver")
	}
	time.Sleep(500 * time.Millisecond)
	logger.Debug("driver started")

	return nil
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

	// ear to read commands
	ear, err := listen.NewEar(conf.Listen.Engine, conf.Listen.Config, logger)
	defer func() {
		ear.Close()
		if err != nil {
			logger.ErrorWith("cannot stop listening").Err("error", err).Write()
		}
	}()
	if err != nil {
		logger.FatalWith("cannot make ear").Err("error", err).Write()
	}

	// voice to say about state
	voice, err := speak.NewVoice(conf.Speak.Engine, conf.Speak.Speaker)
	if err != nil {
		logger.FatalWith("cannot make voice").Err("error", err).Write()
	}

	// body and brain to execute commands
	logger.Info("start thinking")
	body := act.NewBody()
	brain := think.NewBrain(conf.Think.Dry, body, logger)
	defer func() {
		err = brain.Stop()
		if err != nil {
			logger.ErrorWith("cannot stop driver").Err("error", err).Write()
		}
	}()

	// eye to handle video
	eye, err := see.NewEye(conf.See.Engine, conf.See.Config, logger)
	defer func() {
		err = eye.Close()
		if err != nil {
			logger.ErrorWith("cannot close eye").Err("error", err).Write()
		}
	}()
	if err != nil {
		logger.FatalWith("cannot make eye").Err("error", err).Write()
	}

	err = start(logger, body, eye)
	if err != nil {
		logger.FatalWith("cannot start").Err("error", err).Write()
	}

	logger.Info("start")
	for {
		// read command
		text := ear.Listen()
		logger.DebugWith("text heared").String("text", text).Write()

		// parse command
		cmd := command.Understand(text)
		if cmd.Action == "" {
			logger.DebugWith("cannot recognize command").String("text", text).Write()
			continue
		}

		err = voice.Say(string(cmd.Action))
		if err != nil {
			logger.ErrorWith("cannot say").Err("error", err).Write()
		}

		// act
		if cmd.Action == command.Halt {
			logger.Info("halt command received")
			// we're doing halt and clean-up in defers
			return
		}
		err = brain.Do(cmd)
		if err != nil {
			logger.ErrorWith("cannot do action").Err("error", err).Write()
			continue
		}
	}
}
