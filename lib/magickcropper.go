package lib

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"
)

const (
	ShapeMaskCircle    string = "circle"
	ShapeMaskRoundrect string = "rounded"
	ShapeDiamond       string = "diamond"
)

type MagickCropper struct {
	*imagick.MagickWand
}

func (mc MagickCropper) GetAspect() float32 {
	return float32(mc.GetImageWidth()) / float32(mc.GetImageHeight())
}

func (mc MagickCropper) SmartCrop(w, h uint) error {
	var err error

	aspect := mc.GetAspect()

	if w > h {
		err = mc.ResizeImage(w, uint(float32(w)/aspect), imagick.FILTER_LANCZOS2_SHARP)
	} else {
		err = mc.ResizeImage(uint(float32(h)*aspect), h, imagick.FILTER_LANCZOS2_SHARP)
	}

	nw := int(mc.GetImageWidth())
	nh := int(mc.GetImageHeight())

	mc.CropImage(w, h, int(math.Abs(float64(nw-int(w)))/2), int(math.Abs(float64(nh-int(h)))/2))

	return err
}

func (mc MagickCropper) ShapeImage(t string, param float64) error {
	w, h := mc.GetImageWidth(), mc.GetImageHeight()
	result := imagick.NewMagickWand()
	canvas := imagick.NewDrawingWand()
	pw := imagick.NewPixelWand()

	pw.SetColor("none")
	result.NewImage(w, h, pw)

	pw.SetColor("white")
	canvas.SetFillColor(pw)

  defaultShape := func(cnv *imagick.DrawingWand, width uint, height uint) {
    cnv.Rectangle(0, 0, float64(width), float64(height))
  }

	switch t {
	case ShapeMaskCircle:
		if w > h {
			canvas.Circle(float64(w)/2, float64(h)/2, float64(w)/2, float64(h-1))
		} else {
			canvas.Circle(float64(w)/2, float64(h)/2, float64(w-1), float64(h)/2)
		}

	case ShapeMaskRoundrect:
    if param != 0 {
  		canvas.RoundRectangle(0, 0, float64(w), float64(h), param, param)
    } else {
      defaultShape(canvas, w, h)
    }

	case ShapeDiamond:
		canvas.Polygon([]imagick.PointInfo{
			{float64(w) / 2, 0},
			{float64(w), float64(h) / 2},
			{float64(w) / 2, float64(h)},
			{0, float64(h) / 2},
		})

	default:
    defaultShape(canvas, w, h)
	}

	result.DrawImage(canvas)
	result.CompositeImage(mc.MagickWand, imagick.COMPOSITE_OP_SRC_IN, true, 0, 0)

	mc.MagickWand.Clear()
	mc.MagickWand.AddImage(result)

	return nil
}
