package model

type AlbumImage struct {
	Album       string `json:"album" db:"album"`
	Image       string `json:"image" db:"image"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}
