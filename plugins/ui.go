package plugins

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type UI struct {
	app fyne.App
	win fyne.Window

	battery    *canvas.Text
	warns      *canvas.Text
	warnsState map[string]bool
}

var _ StateHandler = &UI{}

func NewUI() *UI {
	return &UI{
		app:        app.New(),
		warnsState: make(map[string]bool),
	}
}

func (ui *UI) Connect(pl *Plugins) {
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
	ui.warns = canvas.NewText("", color.Black)
	ui.win.SetContent(
		container.New(layout.NewGridLayout(1), ui.battery, ui.warns),
	)
	ui.win.Show()
	return nil
}

func (ui *UI) SetBattery(val int8) {
	ui.battery.Text = fmt.Sprintf("battery %d%%", val)
	ui.battery.Refresh()
}

func (ui *UI) SetWarning(msg string, state bool) {
	ui.warnsState[msg] = state

	text := ""
	for msg, state := range ui.warnsState {
		if state {
			text += msg + "\n"
		}
	}
	ui.warns.Text = text
	ui.warns.Refresh()
}
