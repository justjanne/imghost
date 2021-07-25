package model

import (
	"errors"
	"time"
)

type Album struct {
	Id          string    `json:"id" db:"id"`
	Owner       string    `json:"owner" db:"owner"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (album Album) VerifyOwner(user User) error {
	if album.Owner != user.Id {
		return errors.New("user does not have ownership over this item")
	}

	return nil
}
