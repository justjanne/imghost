package task

import (
	"context"
	"encoding/json"
	"fmt"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
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

	println("parsed task: " + payload.ImageId)

	if err = processor.env.Repositories.Images.UpdateState(payload.ImageId, repo.StateInProgress); err != nil {
		println("failed to set image state: " + payload.ImageId)
		println(err.Error())
		return
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	sourceFile, err := ioutil.TempFile("", payload.ImageId+"*."+payload.Extension)
	if err != nil {
		println("failed to create temp file: " + payload.ImageId)
		println(err.Error())
		_ = processor.env.Repositories.Images.UpdateState(payload.ImageId, repo.StateError)
		return
	}
	err = processor.env.Storage.DownloadFile(
		ctx,
		processor.env.Configuration.Storage.ConversionBucket,
		payload.ImageId,
		sourceFile)
	if err != nil {
		println("failed to download file: " + sourceFile.Name())
		println(err.Error())
		_ = processor.env.Repositories.Images.UpdateState(payload.ImageId, repo.StateError)
		return
	}
	if err = wand.ReadImage(sourceFile.Name()); err != nil {
		println("failed to read file: " + sourceFile.Name())
		println(err.Error())
		_ = processor.env.Repositories.Images.UpdateState(payload.ImageId, repo.StateError)
		return
	}
	var originalImage imgconv.ImageHandle
	if originalImage, err = imgconv.NewImage(wand); err != nil {
		println("failed to load file: " + sourceFile.Name())
		println(err.Error())
		_ = processor.env.Repositories.Images.UpdateState(payload.ImageId, repo.StateError)
		return err
	}

	err = util.LaunchGoroutines(len(payload.Sizes), func(index int) error {
		outputFile, err := ioutil.TempFile("", payload.ImageId+"*.png")
		if err != nil {
			return err
		}

		size := payload.Sizes[index]
		image := originalImage.CloneImage()
		if err := image.Crop(size.Size); err != nil {
			return err
		}
		if err := image.Resize(size.Size); err != nil {
			return err
		}
		if err := image.Write(payload.Quality, outputFile); err != nil {
			return err
		}
		if err := processor.env.Storage.UploadFile(
			ctx,
			processor.env.Configuration.Storage.ImageBucket,
			fmt.Sprintf("%s%s", payload.ImageId, size.Suffix),
			"image/png",
			outputFile); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		println("failed to convert image file")
		println(err.Error())
		_ = processor.env.Repositories.Images.UpdateState(payload.ImageId, repo.StateError)
		return
	}

	if err = processor.env.Repositories.Images.UpdateState(payload.ImageId, repo.StateDone); err != nil {
		return
	}

	return
}
