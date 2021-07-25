package task

import (
	"context"
	"encoding/json"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"github.com/hibiken/asynq"
	"github.com/justjanne/imgconv"
	"gopkg.in/gographics/imagick.v2/imagick"
	"io/ioutil"
)

type ImageProcessor struct {
	env environment.BackendEnvironment
}

func NewImageProcessor(env environment.BackendEnvironment) *ImageProcessor {
	return &ImageProcessor{
		env,
	}
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload ImageResizePayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	sourceFile, err := ioutil.TempFile("", payload.ImageId)
	if err != nil {
		return
	}
	err = processor.env.Storage.DownloadFile(
		ctx,
		processor.env.Configuration.Storage.ConversionBucket,
		payload.ImageId,
		sourceFile)
	if err != nil {
		return
	}
	if err = wand.ReadImageFile(sourceFile); err != nil {
		return
	}
	var originalImage imgconv.ImageHandle
	if originalImage, err = imgconv.NewImage(wand); err != nil {
		return err
	}

	err = util.LaunchGoroutines(len(payload.Sizes), func(index int) error {
		outputFile, err := ioutil.TempFile("", payload.ImageId)
		if err != nil {
			return err
		}

		size := payload.Sizes[index]
		image := originalImage.CloneImage()
		if err := image.Crop(size); err != nil {
			return err
		}
		if err := image.Resize(size); err != nil {
			return err
		}
		if err := image.Write(payload.Quality, outputFile); err != nil {
			return err
		}
		if err := processor.env.Storage.UploadFile(
			ctx,
			processor.env.Configuration.Storage.ImageBucket,
			payload.ImageId,
			outputFile); err != nil {
			return err
		}
		return nil
	})

	return
}
