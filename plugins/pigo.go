package plugins

import (
	"fmt"
	"image"
	"time"

	_ "embed"

	pigo "github.com/esimov/pigo/core"
)

//go:embed facefinder-cascade.bin
var cascade []byte

type PiGo struct {
	classifier *pigo.Pigo
	last       time.Time
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

func (PiGo) Connect(pl *Plugins) {
}

func (pg *PiGo) Detect(img image.Image) []image.Rectangle {
	now := time.Now()
	if now.Sub(pg.last) < 2*time.Second {
		return nil
	}
	pg.last = now

	pixels := pigo.RgbToGrayscale(img)
	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   480,
			Cols:   360,
			Dim:    360,
		},
	}
	dets := pg.classifier.RunCascade(cParams, 0)
	dets = pg.classifier.ClusterDetections(dets, 0.2)

	res := make([]image.Rectangle, len(dets))
	for i, det := range dets {
		res[i] = image.Rect(
			det.Col, det.Row,
			det.Col+det.Scale, det.Row+det.Scale,
		)
	}
	return res
}
