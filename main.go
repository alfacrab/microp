package main

import (
	//"fmt"
	L "github.com/absinsekt/mobile-icon-cropper/lib"
	//"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	conf := L.ConfigProvider{}
	conf.Initialize("config.yaml")
	//flow := make(chan L.MagickCropper)
	//batch := make(chan bool)

	// TODO move to yaml config
	//SIZES := [][2]uint{
	//	{32, 32},
	//	{64, 64},
	//	{128, 128},
	//	{128, 64},
	//	{64, 128},
	//}

	//imagick.Initialize()
	//defer imagick.Terminate()
	//
	//mw := imagick.NewMagickWand()
	//
	// TODO move to cli params
	//if err := mw.ReadImage("test.jpg"); err != nil {
	//	panic(err)
	//}

	//for _, pair := range SIZES {
	//	go func(p [2]uint) {
	//		crp := <-flow
	//
	//		crp.SmartCrop(p[0], p[1])
	//
	//		// TODO move to yaml config
	//		crp.ShapeImage(SHAPE_MASK_ROUNDRECT, 10)
	//		crp.MagickWand.WriteImage(
	//			fmt.Sprintf("out_%dx%d.png", p[0], p[1]))
	//
	//		batch<-true
	//	}(pair)
	//
	//	flow <- MagickCropper{mw.Clone()}
	//}
	//
	//for s := 0; s < len(SIZES); s++ {
	//	//fmt.Println("ololo!")
	//	<-batch
	//}
}
