package model

import (
	"errors"
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/storage"
	"time"
)

type Image struct {
	Id           string    `json:"id" db:"id"`
	Owner        string    `json:"owner" db:"owner"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	OriginalName string    `json:"original_name" db:"original_name"`
	MimeType     string    `json:"mime_type" db:"type"`
	State        string    `json:"state" db:"state"`
	Url          string    `json:"url"`
}

func (image Image) VerifyOwner(user User) error {
	if image.Owner != user.Id {
		return errors.New("user does not have ownership over this item")
	}

	return nil
}

func (image *Image) LoadUrl(storage storage.Storage, config configuration.StorageConfiguration) (err error) {
	image.Url = storage.UrlFor(config.ImageBucket, image.Id).String()
	return
}
