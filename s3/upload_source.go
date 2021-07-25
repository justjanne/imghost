package s3

import (
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"github.com/justjanne/imgconv"
)

// TODO: Implement
func UploadSource(env environment.Environment, imageId string, source string) error {
	return nil
}

func DownloadSource(env environment.Environment, imageId string) (string, error) {
	return "", nil
}

func UploadImage(env environment.Environment, imageId string, format imgconv.Size, source string) error {
	return nil
}
