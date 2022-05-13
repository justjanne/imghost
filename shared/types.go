package shared

import (
	"time"
)

type Image struct {
	Id           string `json:"id"`
	Title        string
	Description  string
	CreatedAt    time.Time
	OriginalName string
	MimeType     string `json:"mimeType"`
}

type Result struct {
	Id      string   `json:"id"`
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}
