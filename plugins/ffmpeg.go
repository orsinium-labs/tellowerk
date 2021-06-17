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
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3]
	return color.RGBA{s[2], s[1], s[0], 0}
}

func (p *RGB) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *RGB) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	s := p.Pix[i : i+3 : i+3]
	s[0] = c1.B
	s[1] = c1.G
	s[2] = c1.R
}

type FFMpeg struct {
	logger *zap.Logger
	driver *tello.Driver
	pigo   *PiGo

	in   io.WriteCloser
	out  io.ReadCloser
	win  *imgshow.Window
	dets []image.Rectangle
}

func NewFFMpeg(driver *tello.Driver) *FFMpeg {
	return &FFMpeg{driver: driver}
}

func (ff *FFMpeg) Connect(pl *Plugins) {
	ff.pigo = pl.PiGo
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

		// detect faces
		if ff.pigo != nil {
			dets := ff.pigo.Detect(&img)
			if dets != nil {
				if len(dets) == 0 {
					ff.logger.Debug("faces detected", zap.Int("count", len(dets)))
				}
				ff.dets = dets
			}
		}
		// draw rectangles for detected faces
		c := color.Black
		for _, det := range ff.dets {
			for x := det.Min.X; x < det.Max.X; x++ {
				img.Set(x, det.Min.Y, c)
				img.Set(x, det.Max.Y, c)
			}
			for y := det.Min.Y; y < det.Max.Y; y++ {
				img.Set(det.Min.X, y, c)
				img.Set(det.Max.Y, y, c)
			}
		}

		err = ff.win.Draw(&img)
		if err != nil {
			ff.logger.Error("cannot draw frame: %v", zap.Error(err))
			continue
		}
	}
}
