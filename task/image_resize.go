package task

import (
	"encoding/json"
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"github.com/hibiken/asynq"
	"github.com/justjanne/imgconv"
)

type ImageResizePayload struct {
	ImageId   string
	Extension string
	Sizes     []configuration.SizeDefinition
	Quality   imgconv.Quality
}

func NewResizeTask(imageId string, extension string, config configuration.FrontendConfiguration) (task *asynq.Task, err error) {
	payload, err := json.Marshal(ImageResizePayload{
		ImageId:   imageId,
		Extension: extension,
		Sizes:     config.Conversion.Sizes,
		Quality:   config.Conversion.Quality,
	})
	if err != nil {
		return
	}
	task = asynq.NewTask(config.Conversion.TaskId, payload)
	return
}
