package see

import (
	"errors"
	"image/color"
	"io"
	"os/exec"
	"strconv"

	"github.com/francoispqt/onelog"
	"gocv.io/x/gocv"
)

// OpenCVEye detects face on video stream
// https://gobot.io/documentation/examples/tello_opencv/
// https://github.com/hybridgroup/gobot/blob/master/examples/tello_opencv.go
type OpenCVEye struct {
	logger *onelog.Logger

	classifier gocv.CascadeClassifier
	color      color.RGBA
	window     *gocv.Window

	ffmpegIn  io.WriteCloser
	ffmpegOut io.ReadCloser

	FrameX    int
	FrameY    int
	FrameSize int
}

// Handle processes a new frame on a screen
func (eye *OpenCVEye) Handle(data interface{}) {
	if eye.ffmpegIn == nil {
		return
	}
	raw := data.([]byte)
	_, err := eye.ffmpegIn.Write(raw)
	if err != nil {
		eye.logger.ErrorWith("cannot pipe data into ffmpeg").Err("error", err).Write()
		return
	}
}

func (eye *OpenCVEye) render() {
	eye.logger.Debug("start video rendering goroutine")
	for {
		out := eye.ffmpegOut
		if out == nil {
			return
		}
		// prepare image matrix
		buf := make([]byte, eye.FrameSize)
		_, err := io.ReadFull(out, buf)
		if err != nil {
			if eye.ffmpegOut == nil {
				return
			}
			eye.logger.ErrorWith("cannot read ffmpeg out").Err("error", err).Write()
			continue
		}
		img, err := gocv.NewMatFromBytes(eye.FrameY, eye.FrameX, gocv.MatTypeCV8UC3, buf)
		if err != nil {
			eye.logger.ErrorWith("cannot make matrix").Err("error", err).Write()
			continue
		}
		if img.Empty() {
			eye.logger.Debug("empty matrix")
			continue
		}
		defer img.Close()

		// detect faces
		rects := eye.classifier.DetectMultiScale(img)
		// draw a rectangle around each face on the original image
		for _, r := range rects {
			gocv.Rectangle(&img, r, eye.color, 3)
		}
		eye.window.IMShow(img)
		eye.window.WaitKey(1)
	}
}

// Close closes stdin stream to mlayer
func (eye *OpenCVEye) Close() (err error) {
	err = eye.classifier.Close()
	if err != nil {
		return err
	}
	err = eye.ffmpegIn.Close()
	if err != nil {
		return err
	}
	eye.ffmpegIn = nil
	err = eye.ffmpegOut.Close()
	if err != nil {
		return err
	}
	eye.ffmpegOut = nil

	return eye.window.Close()
}

// NewOpenCVEye creates OpenCVEye instance to detect face on video
func NewOpenCVEye(logger *onelog.Logger, config Config) (*OpenCVEye, error) {
	eye := &OpenCVEye{
		logger:     logger,
		classifier: gocv.NewCascadeClassifier(),
		color:      color.RGBA{0, 0, 255, 0},
		window:     gocv.NewWindow("tellowerk drone"),
		FrameX:     config.FrameX,
		FrameY:     config.FrameY,
		FrameSize:  config.FrameX * config.FrameY * 3,
	}

	ffmpeg := exec.Command(
		"ffmpeg",
		"-hwaccel", "auto",
		"-hwaccel_device", "opencl",
		"-i", "pipe:0",
		"-pix_fmt", "bgr24",
		"-s", strconv.Itoa(config.FrameX)+"x"+strconv.Itoa(config.FrameY),
		"-f", "rawvideo",
		"pipe:1",
	)
	var err error
	eye.ffmpegIn, err = ffmpeg.StdinPipe()
	if err != nil {
		return eye, err
	}
	eye.ffmpegOut, err = ffmpeg.StdoutPipe()
	if err != nil {
		return eye, err
	}
	err = ffmpeg.Start()
	if err != nil {
		return eye, err
	}

	if !eye.classifier.Load(config.Model) {
		return eye, errors.New("cannot read model")
	}

	go eye.render()
	return eye, nil
}
