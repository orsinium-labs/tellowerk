package plugins

import (
	"image"
	"image/color"
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
	return color.RGBA{s[2], s[1], s[0], 255}
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
