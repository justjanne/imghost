package shared

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

const (
	TypeImageResize = "image:resize"
)

type ImageTaskPayload struct {
	ImageId string
}

func NewImageResizeTask(imageId string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageTaskPayload{ImageId: imageId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeImageResize, payload), nil
}
