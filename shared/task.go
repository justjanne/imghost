package shared

import (
	"encoding/json"
	"fmt"
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

func FormatPayload(typeName string, data []byte) string {
	switch typeName {
	case TypeImageResize:
		return string(data)
	default:
		return fmt.Sprintf("unknown type %s", typeName)
	}
}

func FormatResult(typeName string, data []byte) string {
	switch typeName {
	case TypeImageResize:
		return string(data)
	default:
		return fmt.Sprintf("unknown type %s", typeName)
	}
}
