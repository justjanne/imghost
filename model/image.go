package model

import (
	"errors"
	"time"
)

type Image struct {
	Id           string    `json:"id"`
	Owner        string    `json:"owner"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
}

func (image Image) VerifyOwner(user User) error {
	if image.Owner != user.Id {
		return errors.New("user does not have ownership over this item")
	}

	return nil
}
