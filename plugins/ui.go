package plugins

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

type UI struct {
	app fyne.App
	win fyne.Window

	battery    *canvas.Text
	warns      *canvas.Text
	video      *canvas.Image
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
	// ui.app.Settings().SetTheme(theme.LightTheme())
	ui.win = ui.app.NewWindow("tellowerk")

	ui.battery = canvas.NewText("battery", theme.ForegroundColor())
	ui.warns = canvas.NewText("", theme.ForegroundColor())
	ui.video = canvas.NewImageFromImage(
		image.NewRGBA(image.Rect(0, 0, frameX, frameY)),
	)
	ui.video.SetMinSize(fyne.NewSize(frameX, frameY))
	content := container.NewHBox(
		container.New(layout.NewGridLayout(1), ui.battery, ui.warns),
		ui.video,
	)
	ui.win.SetContent(content)
	ui.win.Show()
	return nil
}

func (ui *UI) SetBattery(val int8) {
	if ui.battery == nil {
		return
	}
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

func (ui *UI) SetFrame(img *RGB) {
	ui.video.File = ""
	ui.video.Image = img
	ui.video.Refresh()
}
