package plugins

import (
	"fmt"
	"image"
	"image/color"
	"time"

	_ "embed"

	pigo "github.com/esimov/pigo/core"
	"go.uber.org/zap"
)

//go:embed facefinder-cascade.bin
var cascade []byte

type PiGo struct {
	logger    *zap.Logger
	targeting *Targeting

	classifier *pigo.Pigo
	last       time.Time
	dets       []image.Rectangle
}

func NewPiGo() *PiGo {
	return &PiGo{last: time.Now()}
}

func (pg *PiGo) Start() error {
	var err error
	pigo := pigo.NewPigo()
	pg.classifier, err = pigo.Unpack(cascade)
	if err != nil {
		return fmt.Errorf("unpack cascade file: %v", err)
	}
	return nil
}

func (PiGo) Stop() error {
	return nil
}

func (pg *PiGo) Connect(pl *Plugins) {
	pg.targeting = &Targeting{c: pl.Controller}
}

func (pg PiGo) Dets() []image.Rectangle {
	return pg.dets
}

func (pg *PiGo) Detect(img image.Image) {
	now := time.Now()
	if now.Sub(pg.last) < 400*time.Millisecond {
		return
	}
	pg.last = now
	go pg.detect(img)
}

func (pg *PiGo) detect(img image.Image) {
	pixels := pigo.RgbToGrayscale(img)
	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   frameY,
			Cols:   frameX,
			Dim:    frameX,
		},
	}
	dets := pg.classifier.RunCascade(cParams, 0)
	dets = pg.classifier.ClusterDetections(dets, 0.2)

	res := make([]image.Rectangle, len(dets))
	for i, det := range dets {
		rad := det.Scale / 2
		res[i] = image.Rect(
			det.Col-rad, det.Row-rad,
			det.Col+rad, det.Row+rad,
		)
	}
	if len(res) > 0 {
		pg.logger.Debug("faces detected", zap.Int("count", len(res)))
	}

	err := pg.targeting.Target(res)
	if err != nil {
		pg.logger.Error("target to face", zap.Error(err))
	}

	pg.dets = res
}

func (pg *PiGo) Draw(img *RGB) {
	c := color.RGBA{255, 0, 0, 255}
	for _, det := range pg.Dets() {
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
}
