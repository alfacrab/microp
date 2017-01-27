package main

import (
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"math"
)

const (
	SHAPE_MASK_CIRCLE uint8 = iota
	SHAPE_MASK_ROUNDRECT
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
		err = mc.ResizeImage(w, uint(float32(w)/aspect), imagick.FILTER_LANCZOS2_SHARP, 1)
	} else {
		err = mc.ResizeImage(uint(float32(h)*aspect), h, imagick.FILTER_LANCZOS2_SHARP, 1)
	}

	nw := int(mc.GetImageWidth())
	nh := int(mc.GetImageHeight())

	mc.CropImage(w, h, int(math.Abs(float64(nw-int(w)))/2), int(math.Abs(float64(nh-int(h)))/2))

	return err
}

func (mc MagickCropper) ShapeImage(t uint8, param float64) error {
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
			canvas.Circle(float64(w)/2, float64(h)/2, float64(w)/2, float64(h))
		} else {
			canvas.Circle(float64(w)/2, float64(h)/2, float64(w), float64(h)/2)
		}

	case SHAPE_MASK_ROUNDRECT:
		canvas.RoundRectangle(0, 0, float64(w), float64(h), param, param)
	}

	result.DrawImage(canvas)
	result.CompositeImage(mc.MagickWand, imagick.COMPOSITE_OP_SRC_IN, 0, 0)

	mc.MagickWand.Clear()
	mc.MagickWand.AddImage(result)

	return nil
}

func main() {
	// TODO move to yaml config
	SIZES := [][2]uint{
		{32, 32},
		{64, 64},
		{128, 128},
		{128, 64},
		{64, 128},
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	// TODO move to cli params
	if err := mw.ReadImage("test.jpg"); err != nil {
		panic(err)
	}

	// TODO try goroutines
	for _, pair := range SIZES {
		crp := MagickCropper{mw.Clone()}

		crp.SmartCrop(pair[0], pair[1])

		// TODO move to yaml config
		crp.ShapeImage(SHAPE_MASK_ROUNDRECT, 10)

		crp.MagickWand.WriteImage(
			fmt.Sprintf("out_%dx%d.png", pair[0], pair[1]))
	}
}
