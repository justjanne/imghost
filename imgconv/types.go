package imgconv

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

const (
	ImageFitCover   = "cover"
	ImageFitContain = "contain"
)

type Size struct {
	Width  uint   `json:"width" yaml:"width"`
	Height uint   `json:"height" yaml:"height"`
	Format string `json:"format" yaml:"format"`
}

type Quality struct {
	CompressionQuality uint      `json:"compressionQuality" yaml:"compressionQuality"`
	SamplingFactors    []float64 `json:"samplingFactors" yaml:"samplingFactors"`
}

type ColorProfile struct {
	data   []byte
	format string
}

type ImageHandle struct {
	wand     *imagick.MagickWand
	depth    uint
	profiles []ColorProfile
}
