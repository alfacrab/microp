package lib

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"
)

const (
	SHAPE_MASK_CIRCLE    string = "circle"
	SHAPE_MASK_ROUNDRECT string = "rounded"
	SHAPE_DIAMOND        string = "diamond"
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

	switch t {
	case SHAPE_MASK_CIRCLE:
		if w > h {
			canvas.Circle(float64(w)/2, float64(h)/2, float64(w)/2, float64(h-1))
		} else {
			canvas.Circle(float64(w)/2, float64(h)/2, float64(w-1), float64(h)/2)
		}

	case SHAPE_MASK_ROUNDRECT:
		canvas.RoundRectangle(0, 0, float64(w), float64(h), param, param)

	case SHAPE_DIAMOND:
		canvas.Polygon([]imagick.PointInfo{
			{float64(w) / 2, 0},
			{float64(w), float64(h) / 2},
			{float64(w) / 2, float64(h)},
			{0, float64(h) / 2},
		})

	default:
		canvas.Rectangle(0, 0, float64(w), float64(h))
	}

	result.DrawImage(canvas)
	result.CompositeImage(mc.MagickWand, imagick.COMPOSITE_OP_SRC_IN, true, 0, 0)

	mc.MagickWand.Clear()
	mc.MagickWand.AddImage(result)

	return nil
}
