package model

type ImageInfo struct {
	Image Image  `json:"image"`
	State string `json:"state"`
	Url   string `json:"url"`
}
