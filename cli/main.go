package main

import (
	"encoding/json"
	"flag"
	"git.kuschku.de/justjanne/imghost/imgconv"
	"gopkg.in/gographics/imagick.v2/imagick"
	"io/fs"
	"io/ioutil"
	"os"
)

type arguments struct {
	Width          *uint
	Height         *uint
	Fit            *string
	Quality        *uint
	Source         string
	Target         string
	ExportMetadata *string
}

var args = arguments{
	Width: flag.Uint(
		"width",
		0,
		"Desired width of the image",
	),
	Height: flag.Uint(
		"height",
		0,
		"Desired height of the image",
	),
	Fit: flag.String(
		"fit",
		"contain",
		"Desired fit format for image. Allowed are cover and contain.",
	),
	Quality: flag.Uint(
		"quality",
		90,
		"Desired quality of output image",
	),
	ExportMetadata: flag.String(
		"export-metadata",
		"",
		"Export metadata as json",
	),
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	imagick.Initialize()
	defer imagick.Terminate()

	source := flag.Arg(0)
	target := flag.Arg(1)

	data, err := convert(source, target, imgconv.Quality{
		CompressionQuality: *args.Quality,
		SamplingFactors:    []float64{1.0, 1.0, 1.0, 1.0},
	}, imgconv.Size{
		Width:  *args.Width,
		Height: *args.Height,
		Format: *args.Fit,
	})
	if err != nil {
		panic(err)
	}

	if *args.ExportMetadata != "" {
		marshalled, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(*args.ExportMetadata, marshalled, fs.FileMode(644)); err != nil {
			panic(err)
		}
	}
}

func convert(source string, target string, quality imgconv.Quality, size imgconv.Size) (*imgconv.Metadata, error) {
	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	var err error
	if err = wand.ReadImage(source); err != nil {
		return nil, err
	}
	var image imgconv.ImageHandle
	if image, err = imgconv.NewImage(wand); err != nil {
		return nil, err
	}
	data := image.ParseMetadata()
	if err := image.Crop(size); err != nil {
		return nil, err
	}
	if err := image.Resize(size); err != nil {
		return nil, err
	}
	if err := image.Write(quality, target); err != nil {
		return nil, err
	}
	return &data, nil
}
