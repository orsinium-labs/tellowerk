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

type FFMpeg struct {
	logger    *zap.Logger
	driver    *tello.Driver
	pigo      *PiGo
	targeting *Targeting

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
	ff.targeting = &Targeting{c: pl.Controller}
}

func (ff *FFMpeg) Start() error {
	ffmpeg := exec.Command(
		"ffmpeg",
		"-hwaccel", "auto",
		"-hwaccel_device", "opencl",
		"-i", "pipe:0",
		"-pix_fmt", "bgr24",
		"-s", fmt.Sprintf("%dx%d", frameX, frameY),
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
	c.Width = frameX
	c.Height = frameY
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
		buf := make([]byte, frameX*frameY*3)
		_, err = io.ReadFull(ff.out, buf)
		if err != nil {
			ff.logger.Error("cannot read ffmpeg stdout", zap.Error(err))
			continue
		}
		img := RGB{
			Pix:    []uint8(buf),
			Stride: frameX * 3,
			Rect:   image.Rect(0, 0, frameX, frameY),
		}

		// detect faces
		if ff.pigo != nil {
			dets := ff.pigo.Detect(&img)
			if dets != nil {
				if len(dets) != 0 {
					ff.logger.Debug("faces detected", zap.Int("count", len(dets)))
					err = ff.targeting.Target(dets)
					if err != nil {
						ff.logger.Error("cannot target to face", zap.Error(err))
					}
				}
				ff.dets = dets
			}
		}
		// draw rectangles for detected faces
		c := color.RGBA{255, 0, 0, 255}
		for _, det := range ff.dets {
			for x := det.Min.X; x < det.Max.X; x++ {
				img.Set(x, det.Min.Y, c)
				img.Set(x, det.Min.Y+1, c)
				img.Set(x, det.Max.Y, c)
				img.Set(x, det.Max.Y+1, c)
			}
			for y := det.Min.Y; y < det.Max.Y; y++ {
				img.Set(det.Min.X, y, c)
				img.Set(det.Min.X+1, y, c)
				img.Set(det.Max.X, y, c)
				img.Set(det.Max.X+1, y, c)
			}
		}

		err = ff.win.Draw(&img)
		if err != nil {
			ff.logger.Error("cannot draw frame: %v", zap.Error(err))
			continue
		}
	}
}
