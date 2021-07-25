package configuration

type StorageConfiguration struct {
	Endpoint         string `json:"endpoint" yaml:"endpoint"`
	Secure           bool   `json:"secure" yaml:"secure"`
	ImageBucket      string `json:"image_bucket" yaml:"image-bucket"`
	ConversionBucket string `json:"conversion_bucket" yaml:"conversion-bucket"`
	AccessKey        string `json:"access_key" yaml:"access-key"`
	SecretKey        string `json:"secret_key" yaml:"secret-key"`
}
