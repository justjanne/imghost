package shared

import (
	"git.kuschku.de/justjanne/imghost/imgconv"
	"github.com/hibiken/asynq"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type SizeDefinition struct {
	Size   imgconv.Size `yaml:"size"`
	Suffix string       `yaml:"suffix"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type DatabaseConfig struct {
	Format string `yaml:"format"`
	Url    string `yaml:"url"`
}

type Config struct {
	Sizes         []SizeDefinition `yaml:"sizes"`
	Quality       imgconv.Quality  `yaml:"quality"`
	SourceFolder  string           `yaml:"sourceFolder"`
	TargetFolder  string           `yaml:"targetFolder"`
	Redis         RedisConfig      `yaml:"redis"`
	Database      DatabaseConfig   `yaml:"database"`
	Concurrency   int              `yaml:"concurrency"`
	UploadTimeout string           `yaml:"uploadTimeout"`
	BaseUrl       string           `yaml:"baseUrl"`
}

func LoadConfigFromFile(file *os.File) Config {
	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalf("error loading config, %s", err.Error())
	}
	return config
}

func (config Config) UploadTimeoutDuration() time.Duration {
	duration, err := time.ParseDuration(config.UploadTimeout)
	if err != nil {
		log.Fatalf("error loading config: error parsing upload timeout, %s", err.Error())
	}
	return duration
}

func (config Config) AsynqOpts() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.Database,
	}
}
