package task

import (
	"context"
	"encoding/json"
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/s3"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"github.com/hibiken/asynq"
	"github.com/justjanne/imgconv"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type ImageResizePayload struct {
	ImageId string
	Sizes   []imgconv.Size
	Quality imgconv.Quality
}

func NewResizeTask(imageId string, config configuration.Configuration) (task *asynq.Task, err error) {
	payload, err := json.Marshal(ImageResizePayload{
		ImageId: imageId,
		Sizes:   config.Conversion.Sizes,
		Quality: config.Conversion.Quality,
	})
	if err != nil {
		return
	}
	task = asynq.NewTask(config.Conversion.ResizeTaskId, payload)
	return
}

func HandleImageResizeTask(ctx context.Context, task *asynq.Task) (err error) {
	// TODO: Handle environment for tasks
	env := environment.Environment{}

	var payload ImageResizePayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	file, err := s3.DownloadSource(env, payload.ImageId)
	if err != nil {
		return
	}
	if err = wand.ReadImage(file); err != nil {
		return
	}
	var originalImage imgconv.ImageHandle
	if originalImage, err = imgconv.NewImage(wand); err != nil {
		return err
	}

	err = util.LaunchGoroutines(len(payload.Sizes), func(index int) error {
		size := payload.Sizes[index]
		// TODO: Allocate temp file
		tmpFile := ""
		image := originalImage.CloneImage()
		if err := image.Crop(size); err != nil {
			return err
		}
		if err := image.Resize(size); err != nil {
			return err
		}
		if err := image.Write(payload.Quality, tmpFile); err != nil {
			return err
		}
		if err := s3.UploadImage(env, payload.ImageId, size, tmpFile); err != nil {
			return err
		}
		return nil
	})

	return
}
