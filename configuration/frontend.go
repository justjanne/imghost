package configuration

type FrontendConfiguration struct {
	Queue      QueueConfiguration      `json:"queue" yaml:"queue"`
	Database   DatabaseConfiguration   `json:"database" yaml:"database"`
	Redis      RedisConfiguration      `json:"redis" yaml:"redis"`
	Conversion ConversionConfiguration `json:"conversion" yaml:"conversion"`
	Storage    StorageConfiguration    `json:"storage" yaml:"storage"`
	Auth       AuthConfiguration       `json:"auth" yaml:"auth"`
}
