package imgconv

import (
	"gopkg.in/gographics/imagick.v2/imagick"
	"log"
	"math"
	"os"
	"strings"
)

func NewImage(wand *imagick.MagickWand) (ImageHandle, error) {
	meta := ImageHandle{
		wand:  wand,
		depth: wand.GetImageDepth(),
	}

	if err := wand.AutoOrientImage(); err != nil {
		return meta, err
	}

	if len(wand.GetImageProfiles("i*")) == 0 {
		if err := wand.ProfileImage("icc", ProfileSRGB); err != nil {
			return meta, err
		}
	}

	for _, name := range wand.GetImageProfiles("*") {
		meta.profiles = append(meta.profiles, ColorProfile{
			data:   []byte(wand.GetImageProfile(name)),
			format: name,
		})
	}
	if meta.depth < 16 {
		if err := wand.SetImageDepth(16); err != nil {
			return meta, err
		}
	}
	if err := wand.ProfileImage("icc", ProfileACESLinear); err != nil {
		return meta, err
	}
	return meta, nil
}

func (image *ImageHandle) CloneImage() ImageHandle {
	return ImageHandle{
		image.wand.Clone(),
		image.depth,
		image.profiles,
	}
}

func (image *ImageHandle) ParseMetadata() Metadata {
	return parseMetadata(image.wand)
}

func (image *ImageHandle) SanitizeMetadata() error {
	var profiles []ColorProfile
	for _, profile := range image.profiles {
		if !strings.EqualFold("exif", profile.format) {
			profiles = append(profiles, profile)
		}
	}
	image.profiles = profiles
	image.wand.RemoveImageProfile("exif")

	if err := image.wand.SetOption("png:include-chunk", "bKGD,cHRM,iCCP"); err != nil {
		return err
	}
	if err := image.wand.SetOption("png:exclude-chunk", "EXIF,iTXt,tEXt,zTXt,date"); err != nil {
		return err
	}
	for _, key := range image.wand.GetImageProperties("png:*") {
		if err := image.wand.DeleteImageProperty(key); err != nil {
			return err
		}
	}

	return nil
}

func (image *ImageHandle) Crop(size Size) error {
	if size.Width == 0 || size.Height == 0 || size.Format != ImageFitCover {
		return nil
	}

	currentWidth := image.wand.GetImageWidth()
	currentHeight := image.wand.GetImageHeight()

	currentAspectRatio := float64(currentWidth) / float64(currentHeight)
	desiredAspectRatio := float64(size.Width) / float64(size.Height)

	if currentAspectRatio == desiredAspectRatio {
		return nil
	}

	var desiredWidth, desiredHeight uint
	if desiredAspectRatio > currentAspectRatio {
		desiredWidth = currentWidth
		desiredHeight = uint(math.Round(float64(currentWidth) / desiredAspectRatio))
	} else {
		desiredHeight = currentHeight
		desiredWidth = uint(math.Round(desiredAspectRatio * float64(currentHeight)))
	}

	offsetLeft := int((currentWidth - desiredWidth) / 2.0)
	offsetTop := int((currentHeight - desiredHeight) / 2.0)

	if err := image.wand.CropImage(desiredWidth, desiredHeight, offsetLeft, offsetTop); err != nil {
		return err
	}

	return nil
}

func determineDesiredSize(width uint, height uint, size Size) (uint, uint) {
	currentAspectRatio := float64(width) / float64(height)

	var desiredWidth, desiredHeight uint
	if size.Height != 0 && size.Width != 0 {
		if size.Format == ImageFitCover {
			var desiredAspectRatio = float64(size.Width) / float64(size.Height)
			var croppedWidth, croppedHeight uint
			if desiredAspectRatio > currentAspectRatio {
				croppedWidth = width
				croppedHeight = uint(math.Round(float64(width) / desiredAspectRatio))
			} else {
				croppedHeight = height
				croppedWidth = uint(math.Round(desiredAspectRatio * float64(height)))
			}

			desiredHeight = uint(math.Min(float64(size.Height), float64(croppedHeight)))
			desiredWidth = uint(math.Min(float64(size.Width), float64(croppedWidth)))
		} else if currentAspectRatio > 1 {
			desiredWidth = uint(math.Min(float64(size.Width), float64(width)))
			desiredHeight = uint(math.Round(float64(desiredWidth) / currentAspectRatio))
		} else {
			desiredHeight = uint(math.Min(float64(size.Height), float64(height)))
			desiredWidth = uint(math.Round(currentAspectRatio * float64(desiredHeight)))
		}
	} else if size.Height != 0 {
		desiredHeight = uint(math.Min(float64(size.Height), float64(height)))
		desiredWidth = uint(math.Round(currentAspectRatio * float64(desiredHeight)))
	} else if size.Width != 0 {
		desiredWidth = uint(math.Min(float64(size.Width), float64(width)))
		desiredHeight = uint(math.Round(float64(desiredWidth) / currentAspectRatio))
	} else {
		desiredWidth = width
		desiredHeight = height
	}

	return desiredWidth, desiredHeight
}

func (image *ImageHandle) Resize(size Size) error {
	if size.Width == 0 && size.Height == 0 {
		return nil
	}

	currentWidth := image.wand.GetImageWidth()
	currentHeight := image.wand.GetImageHeight()

	desiredWidth, desiredHeight := determineDesiredSize(currentWidth, currentHeight, size)

	if desiredWidth != currentWidth || desiredHeight != currentHeight {
		if err := image.wand.ResizeImage(desiredWidth, desiredHeight, imagick.FILTER_LANCZOS, 1); err != nil {
			return err
		}
	}

	return nil
}

func (image *ImageHandle) prepareWrite(quality Quality) error {
	for _, profile := range image.profiles {
		if err := image.wand.ProfileImage(profile.format, profile.data); err != nil {
			return err
		}
	}
	if err := image.wand.SetImageDepth(image.depth); err != nil {
		return err
	}

	if quality.CompressionQuality != 0 {
		if err := image.wand.SetImageCompressionQuality(quality.CompressionQuality); err != nil {
			return err
		}
	}

	if len(quality.SamplingFactors) != 0 {
		if err := image.wand.SetSamplingFactors(quality.SamplingFactors); err != nil {
			return err
		}
	}

	log.Printf("done preparing image for writing")

	return nil
}

func (image *ImageHandle) Write(quality Quality, target string) error {
	if err := image.prepareWrite(quality); err != nil {
		return err
	}
	if err := image.wand.WriteImage(target); err != nil {
		return err
	}
	return nil
}

func (image *ImageHandle) WriteImageFile(quality Quality, target *os.File) error {
	if err := image.prepareWrite(quality); err != nil {
		return err
	}
	if err := image.wand.WriteImageFile(target); err != nil {
		return err
	}
	return nil
}
