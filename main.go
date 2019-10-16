package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/orsinium/tellowerk/listen"
	"github.com/orsinium/tellowerk/speak"
)

type configListen struct {
	listen.PocketSphinxConfig
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
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}
	ears, err := listen.NewEars(conf.Listen.Engine, conf.Listen.PocketSphinxConfig)
	if err != nil {
		log.Fatal(err)
	}
	voice := speak.NewVoice(conf.Speak.Engine, conf.Speak.Speaker)
	for {
		text := ears.Listen()
		fmt.Println(text)
		voice.Say(text)
	}
}
