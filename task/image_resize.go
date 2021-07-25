package task

import (
	"context"
	"encoding/json"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"github.com/hibiken/asynq"
	"github.com/justjanne/imgconv"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const TypeImageResize = "image:Resize"

type ImageResizePayload struct {
	ImageId string
	Sizes   []imgconv.Size
	Quality imgconv.Quality
}

func NewResizeTask(imageId string, sizes []imgconv.Size, quality imgconv.Quality) (task *asynq.Task, err error) {
	payload, err := json.Marshal(ImageResizePayload{
		ImageId: imageId,
		Sizes:   sizes,
		Quality: quality,
	})
	if err != nil {
		return
	}
	task = asynq.NewTask(TypeImageResize, payload)
	return
}

func HandleImageResizeTask(ctx context.Context, task *asynq.Task) (err error) {
	var payload ImageResizePayload
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return
	}

	wand := imagick.NewMagickWand()
	defer wand.Destroy()

	tmpFile := ""
	if err = wand.ReadImage(tmpFile); err != nil {
		return
	}
	var originalImage imgconv.ImageHandle
	if originalImage, err = imgconv.NewImage(wand); err != nil {
		return err
	}

	err = util.LaunchGoroutines(len(payload.Sizes), func(index int) error {
		size := payload.Sizes[index]
		tmpTargetFile := ""
		image := originalImage.CloneImage()
		if err := image.Crop(size); err != nil {
			return err
		}
		if err := image.Resize(size); err != nil {
			return err
		}
		if err := image.Write(payload.Quality, tmpTargetFile); err != nil {
			return err
		}
		return nil
	})

	return
}
