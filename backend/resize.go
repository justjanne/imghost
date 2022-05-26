package main

import (
	"fmt"
	"git.kuschku.de/justjanne/imghost/imgconv"
	"gopkg.in/gographics/imagick.v2/imagick"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ResizeImage(env ProcessingEnvironment, imageId string) []error {
	var err error

	log.Printf("creating magick wand for %s", imageId)
	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	startRead := time.Now().UTC()
	log.Printf("reading image for %s", imageId)
	if err = wand.ReadImage(filepath.Join(env.Config.SourceFolder, imageId)); err != nil {
		return []error{err}
	}
	log.Printf("importing image for %s", imageId)
	var originalImage imgconv.ImageHandle
	if originalImage, err = imgconv.NewImage(wand); err != nil {
		return []error{err}
	}
	trackTimeSince(imageProcessDurationRead, startRead)

	log.Printf("launching resize goroutines for %s", imageId)
	return runMany(len(env.Config.Sizes), func(index int) error {
		definition := env.Config.Sizes[index]
		path := filepath.Join(env.Config.TargetFolder, fmt.Sprintf("%s%s", imageId, definition.Suffix))
		startClone := time.Now().UTC()
		log.Printf("cloning image for %s in %v", imageId, definition)
		image := originalImage.CloneImage()
		startCrop := trackTimeSince(imageProcessDurationClone, startClone)
		log.Printf("cropping image for %s in %v", imageId, definition)
		if err := image.Crop(definition.Size); err != nil {
			return err
		}
		startResize := trackTimeSince(imageProcessDurationCrop, startCrop)
		log.Printf("resizing image for %s in %v", imageId, definition)
		if err := image.Resize(definition.Size); err != nil {
			return err
		}
		startWrite := trackTimeSince(imageProcessDurationResize, startResize)
		log.Printf("opening image for %s in %v", imageId, definition)
		target, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		log.Printf("writing image for %s in %v", imageId, definition)
		if err := image.WriteImageFile(env.Config.Quality, target); err != nil {
			return err
		}
		log.Printf("tracking time for %s in %v", imageId, definition)
		trackTimeSince(imageProcessDurationWrite, startWrite)
		log.Printf("done with image for %s in %v", imageId, definition)
		return nil
	})
}
