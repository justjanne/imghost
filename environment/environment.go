package environment

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/justjanne/imgconv"
	"os"
)

type Environment struct {
	Queue      *asynq.Client
	Database   *sqlx.DB
	RolePrefix string
	Sizes      []imgconv.Size
	Quality    imgconv.Quality
}

func NewEnvironment(
	redisAddress string,
	dbType string,
	dbUrl string,
	rolePrefix string,
	sizes []imgconv.Size,
	quality imgconv.Quality,
) (env Environment, err error) {
	env.Queue = asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddress,
	})
	env.Database, err = sqlx.Open(dbType, dbUrl)
	env.RolePrefix = rolePrefix
	env.Sizes = sizes
	env.Quality = quality
	return
}

func InitializeEnvironment() (env Environment, err error) {
	var sizes []imgconv.Size
	if err = json.Unmarshal([]byte(os.Getenv("CONFIG_SIZES")), &sizes); err != nil {
		return
	}
	var quality imgconv.Quality
	if err = json.Unmarshal([]byte(os.Getenv("CONFIG_QUALITY")), &quality); err != nil {
		return
	}
	rolePrefix := os.Getenv("CONFIG_ROLE_PREFIX")

	redisAddress := os.Getenv("REDIS_ADDRESS")
	dbType := os.Getenv("DB_TYPE")
	dbUrl := os.Getenv("DB_URL")
	return NewEnvironment(redisAddress, dbType, dbUrl, rolePrefix, sizes, quality)
}

func (env Environment) Destroy() error {
	if err := env.Queue.Close(); err != nil {
		return err
	}
	if err := env.Database.Close(); err != nil {
		return err
	}
	return nil
}
