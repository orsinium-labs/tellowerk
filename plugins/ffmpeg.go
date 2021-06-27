package plugins

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/orsinium-labs/imgshow"
	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type FFMpeg struct {
	logger    *zap.Logger
	driver    *tello.Driver
	pigo      *PiGo
	targeting *Targeting
	state     *State

	in   io.WriteCloser
	out  io.ReadCloser
	win  *imgshow.Window
	dets []image.Rectangle

	Show bool // if the video should be rendered using imgshow
}

func NewFFMpeg(driver *tello.Driver) *FFMpeg {
	return &FFMpeg{driver: driver}
}

func (ff *FFMpeg) Connect(pl *Plugins) {
	ff.pigo = pl.PiGo
	ff.logger = pl.Logger
	ff.state = pl.State
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

	if ff.Show {
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
	}

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
	if ff.win != nil {
		ff.win.Destroy()
	}
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
	i := 0
	for {
		i = (i + 1) % 3
		if i == 0 {
			continue
		}
		if ff.in == nil {
			return
		}
		err := ff.handleFrame()
		if err != nil {
			ff.logger.Error("cannot process video frame", zap.Error(err))
			continue
		}
	}
}

func (ff *FFMpeg) handleFrame() error {
	var err error
	// read raw frame
	buf := make([]byte, frameX*frameY*3)
	_, err = io.ReadFull(ff.out, buf)
	if err != nil {
		return fmt.Errorf("read ffmpeg stdout: %v", err)
	}
	img := RGB{
		Pix:    []uint8(buf),
		Stride: frameX * 3,
		Rect:   image.Rect(0, 0, frameX, frameY),
	}

	// detect faces
	if ff.pigo != nil && ff.state.FaceCapture() {
		dets := ff.pigo.Detect(&img)
		if dets != nil {
			if len(dets) != 0 {
				ff.logger.Debug("faces detected", zap.Int("count", len(dets)))
			}
			err = ff.targeting.Target(dets)
			if err != nil {
				return fmt.Errorf("target to face: %v", err)
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

	// take a photo
	if ff.state.TakePhoto() {
		ff.state.SetTakePhoto(false)
		fname := fmt.Sprintf("tello-%s.jpg", time.Now().Format("2006-01-02_15-04-05"))
		ff.logger.Debug("save the frame", zap.String("name", fname))
		stream, err := os.Create(fname)
		if err != nil {
			return fmt.Errorf("open file: %v", zap.Error(err))
		}
		defer stream.Close()
		err = jpeg.Encode(stream, &img, nil)
		if err != nil {
			return fmt.Errorf("encode jpeg: %v", zap.Error(err))
		}
	}

	// render the frame on the screen
	if ff.win != nil {
		err = ff.win.Draw(&img)
		if err != nil {
			return fmt.Errorf("draw frame: %v", zap.Error(err))
		}
	}
	return nil
}
