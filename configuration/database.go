package configuration

type DatabaseConfiguration struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}
