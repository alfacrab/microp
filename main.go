package main

import (
	"flag"
	"fmt"
	L "github.com/absinsekt/microp/lib"
	"gopkg.in/gographics/imagick.v3/imagick"
	"os"
	"path/filepath"
)

func main() {
	var (
		srcImageFile    string
		configFile      string
		targetDirectory string
		concurrency     uint
	)

	flag.StringVar(&configFile, "f", "config.yaml", "configuration file in yaml format (default: config.yaml)")
	flag.StringVar(&targetDirectory, "d", "out", "target directory (default: out)")
	flag.UintVar(&concurrency, "c", 5, "batch concurrency (default: 5)")
	flag.Parse()

	if concurrency > 10 || concurrency < 1 {
		concurrency = 5
	}

	args := flag.Args()

	if len(args) == 0 {
		notifyError("source file not set")
	} else {
		srcImageFile = args[0]
	}

	conf := L.ConfigProvider{}
	if err := conf.Initialize(configFile); err != nil {
		notifyError(err)
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	if err := mw.ReadImage(srcImageFile); err != nil {
		notifyError(err)
	}

	batch := make(chan L.MagickCropper, concurrency)
	done := make(chan string)

	for _, set := range conf.ConfigData.Sets {
		for _, icon := range set.Icons {
			var iconName string

			if icon.Name == "" {
				iconName = fmt.Sprintf("%s_%dx%d.png", set.Prefix, icon.Width, icon.Height)
			} else {
				iconName = fmt.Sprintf("%s.png", icon.Name)
			}

			td := filepath.Join(targetDirectory, set.Prefix)
			if _, err := os.Stat(td); os.IsNotExist(err) {
				if err := os.MkdirAll(td, 0755); err != nil {
					notifyError(err)
				}
			}

			go func(tf string, icfg L.IconConfig, scfg L.IconSet) {
				crp := <-batch

				if scfg.RemoveAlpha {
					if crp.GetImageAlphaChannel() {
						fmt.Printf("Warning: alpha channel removed from source at `%s`\n", tf)
						crp.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE)
					}
				}

				crp.SmartCrop(icfg.Width, icfg.Height)
				crp.ShapeImage(icfg.Type, float64(icfg.Radius))

				if err := crp.MagickWand.WriteImage(tf); err != nil {
					notifyError(err)
				}

				done <- tf

			}(filepath.Join(td, iconName), icon, set)

			batch <- L.MagickCropper{mw.Clone()}
		}
	}

	for i := 0; i < conf.ConfigData.Length(); i++ {
		fmt.Printf("file ready: %s\n", <-done)
	}
}

func notifyError(err interface{}) {
	fmt.Println(err)
	os.Exit(1)
}
