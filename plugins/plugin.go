package plugins

import (
	"fmt"
	"reflect"

	"github.com/orsinium-labs/tellowerk/controllers"
	"go.uber.org/zap"
)

type Plugin interface {
	Connect(*Plugins)
	Start() error
	Stop() error
}

type Plugins struct {
	Controller controllers.Controller
	Logger     *zap.Logger

	State    *State
	GamePad  *GamePad
	Video    *Video
	MPlayer  *MPlayer
	Recorder *Recorder
	FFMpeg   *FFMpeg
	PiGo     *PiGo
	UI       *UI
}

func (plugins *Plugins) All() []Plugin {
	return []Plugin{
		plugins.State,
		plugins.GamePad,
		plugins.Video,
		plugins.MPlayer,
		plugins.Recorder,
		plugins.FFMpeg,
		plugins.PiGo,
		plugins.UI,
	}
}

func (plugins *Plugins) Run() error {
	for _, pl := range plugins.All() {
		if reflect.ValueOf(pl).IsNil() {
			continue
		}
		pl.Connect(plugins)
	}

	for _, pl := range plugins.All() {
		if reflect.ValueOf(pl).IsNil() {
			continue
		}
		err := pl.Start()
		if err != nil {
			return fmt.Errorf("start plugin: %v", err)
		}
		defer func(pl Plugin) {
			err := pl.Stop()
			if err != nil {
				plugins.Logger.Error("cannot stop plugin", zap.Error(err))
			}
		}(pl)
	}
	plugins.UI.Wait()
	return nil
}
