package plugins

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os/exec"

	"github.com/orsinium-labs/imgshow"
	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type RGB struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func (p *RGB) ColorModel() color.Model {
	return color.RGBAModel
}

func (p *RGB) Bounds() image.Rectangle {
	return p.Rect
}

func (p *RGB) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
	s := p.Pix[i : i+3 : i+3]
	return color.RGBA{s[2], s[1], s[0], 0}
}

type FFMpeg struct {
	logger *zap.Logger
	driver *tello.Driver

	in  io.WriteCloser
	out io.ReadCloser
	win *imgshow.Window
}

func NewFFMpeg(driver *tello.Driver) *FFMpeg {
	return &FFMpeg{driver: driver}
}

func (ff *FFMpeg) Connect(pl *Plugins) {
	ff.logger = pl.Logger
}

func (ff *FFMpeg) Start() error {
	ffmpeg := exec.Command(
		"ffmpeg",
		"-hwaccel", "auto",
		"-hwaccel_device", "opencl",
		"-i", "pipe:0",
		"-pix_fmt", "bgr24",
		"-s", "480x360",
		"-f", "rawvideo",
		"pipe:1",
	)
	var err error
	ff.in, err = ffmpeg.StdinPipe()
	if err != nil {
		return fmt.Errorf("pipe stdin: %v", err)
	}
	ff.out, err = ffmpeg.StdoutPipe()
	if err != nil {
		return fmt.Errorf("pipe stdout: %v", err)
	}
	err = ffmpeg.Start()
	if err != nil {
		return fmt.Errorf("run ffmpeg: %v", err)
	}

	c := imgshow.NewConfig()
	c.Width = 480
	c.Height = 360
	c.Title = "tellowerk"
	ff.win = c.Window()
	err = ff.win.Create()
	if err != nil {
		return fmt.Errorf("create window: %v", err)
	}
	go ff.win.Render()

	go ff.worker()
	err = ff.driver.On(tello.VideoFrameEvent, ff.handle)
	if err != nil {
		return fmt.Errorf("subscribe to video frames: %v", err)
	}

	return nil
}

func (ff *FFMpeg) handle(data interface{}) {
	if ff.in == nil {
		return
	}
	raw := data.([]byte)
	_, err := ff.in.Write(raw)
	if err != nil {
		ff.logger.Error("cannot pipe data into ffmpeg", zap.Error(err))
		return
	}
}
func (ff *FFMpeg) Stop() error {
	var err error
	ff.win.Destroy()
	err = ff.in.Close()
	if err != nil {
		return fmt.Errorf("close ffmpeg stdin: %v", err)
	}
	ff.in = nil
	err = ff.out.Close()
	if err != nil {
		return fmt.Errorf("close ffmpeg stdout: %v", err)
	}
	ff.out = nil
	return nil
}

func (ff *FFMpeg) worker() {
	var err error
	for {
		if ff.in == nil {
			return
		}
		// read raw frame
		buf := make([]byte, 480*360*3)
		_, err = io.ReadFull(ff.out, buf)
		if err != nil {
			ff.logger.Error("cannot read ffmpeg stdout", zap.Error(err))
			continue
		}
		img := RGB{
			Pix:    []uint8(buf),
			Stride: 480 * 3,
			Rect:   image.Rect(0, 0, 480, 360),
		}
		err = ff.win.Draw(&img)
		if err != nil {
			ff.logger.Error("cannot draw frame: %v", zap.Error(err))
			continue
		}
	}
}
