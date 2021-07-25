package configuration

type DatabaseConfiguration struct {
	Type string `json:"type" yaml:"type"`
	Url  string `json:"url" yaml:"url"`
}
