package plugins

import (
	"image"

	"github.com/orsinium-labs/tellowerk/controllers"
)

type Targeting struct {
	c controllers.Controller
}

func (t *Targeting) Target(targets []image.Rectangle) error {
	err := t.c.Hover()
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		return nil
	}

	target := t.best(targets)
	icenter := image.Point{frameX / 2, frameY / 2}

	// rotate ox
	if target.Min.X > icenter.X {
		err = t.c.Rotate(50)
		if err != nil {
			return err
		}
	} else if target.Max.X < icenter.X {
		err = t.c.Rotate(-50)
		if err != nil {
			return err
		}
	}

	// position oy
	if target.Min.Y > icenter.Y {
		err = t.c.OZ(-20)
		if err != nil {
			return err
		}
	} else if target.Max.Y < icenter.Y {
		err = t.c.OZ(20)
		if err != nil {
			return err
		}
	}

	// center := image.Point{
	// 	X: target.Min.X + target.Dx()/2,
	// 	Y: target.Min.Y + target.Dy()/2,
	// }
	return nil
}

func (t *Targeting) best(targets []image.Rectangle) image.Rectangle {
	if len(targets) == 1 {
		return targets[0]
	}
	var best image.Rectangle
	bestA := 0
	for _, target := range targets[1:] {
		s := target.Size()
		area := s.X * s.Y
		if area > bestA {
			best = target
		}
	}
	return best
}
