package model

import (
	"errors"
	"time"
)

type Album struct {
	Id          string    `json:"id"`
	Owner       string    `json:"owner"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (album Album) VerifyOwner(user User) error {
	if album.Owner != user.Id {
		return errors.New("user does not have ownership over this item")
	}

	return nil
}
