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
	CompressionQuality uint      `json:"compression_quality" yaml:"compression-quality"`
	SamplingFactors    []float64 `json:"sampling_factors" yaml:"sampling-factors"`
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
