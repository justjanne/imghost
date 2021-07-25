package configuration

type Configuration struct {
	Queue      QueueConfiguration      `json:"queue"`
	Database   DatabaseConfiguration   `json:"database"`
	Redis      RedisConfiguration      `json:"redis"`
	Conversion ConversionConfiguration `json:"conversion"`
	Auth       AuthConfiguration       `json:"auth"`
}
