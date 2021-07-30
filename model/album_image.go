package model

import (
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/storage"
)

type AlbumImage struct {
	Album       string `json:"album" db:"album"`
	Image       string `json:"image" db:"image"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Url         string `json:"url"`
}

func (image *AlbumImage) LoadUrl(storage storage.Storage, config configuration.StorageConfiguration) (err error) {
	image.Url = storage.UrlFor(config.ImageBucket, image.Image).String()
	return
}
