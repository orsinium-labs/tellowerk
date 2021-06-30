package plugins

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"time"

	"go.uber.org/zap"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type FFMpeg struct {
	logger *zap.Logger
	driver *tello.Driver
	pigo   *PiGo
	state  *State
	ui     *UI

	in  io.WriteCloser
	out io.ReadCloser
}

func NewFFMpeg(driver *tello.Driver) *FFMpeg {
	return &FFMpeg{driver: driver}
}

func (ff *FFMpeg) Connect(pl *Plugins) {
	ff.pigo = pl.PiGo
	ff.logger = pl.Logger
	ff.state = pl.State
	ff.ui = pl.UI
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
		ff.pigo.Detect(&img)
		ff.pigo.Draw(&img)
	}

	// take a photo
	if ff.state.TakePhoto() {
		ff.state.SetTakePhoto(false)
		fname := fmt.Sprintf("tello-%s.jpg", time.Now().Format("2006-01-02_15-04-05"))
		ff.logger.Debug("save the frame", zap.String("name", fname))
		stream, err := os.Create(fname)
		if err != nil {
			return fmt.Errorf("open file: %v", err)
		}
		defer stream.Close()
		err = jpeg.Encode(stream, &img, nil)
		if err != nil {
			return fmt.Errorf("encode jpeg: %v", err)
		}
	}

	// render the frame
	if ff.ui != nil {
		ff.ui.SetFrame(&img)
	}
	return nil
}
