package plugins

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type UI struct {
	state    *State
	oldState State

	app fyne.App
	win fyne.Window

	battery *canvas.Text
}

func NewUI() *UI {
	return &UI{
		app: app.New(),
	}
}

func (ui *UI) Connect(pl *Plugins) {
	ui.state = pl.State
	ui.oldState = *pl.State
}

func (ui *UI) Wait() {
	ui.app.Run()
}

func (ui *UI) Stop() error {
	ui.app.Quit()
	return nil
}

func (ui *UI) Start() error {
	ui.win = ui.app.NewWindow("tellowerk")

	ui.battery = canvas.NewText("battery", color.Black)
	ui.win.SetContent(
		container.New(layout.NewGridLayout(1), ui.battery),
	)
	ui.win.Show()
	go ui.worker()
	return nil
}

func (ui *UI) worker() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		ui.update()
	}
}

func (ui *UI) update() {
	if ui.state.battery != ui.oldState.battery {
		ui.battery.Text = fmt.Sprintf("battery %d%%", ui.state.battery)
		ui.battery.Refresh()
	}
	ui.oldState = *ui.state
}
