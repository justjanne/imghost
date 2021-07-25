package environment

import (
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
	"git.kuschku.de/justjanne/imghost-frontend/storage"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"time"
)

type FrontendEnvironment struct {
	Configuration configuration.FrontendConfiguration
	QueueClient   *asynq.Client
	Database      *sqlx.DB
	Repositories  Repositories
	Storage       storage.Storage
}

func NewFrontendEnvironment(config configuration.FrontendConfiguration) (env FrontendEnvironment, err error) {
	env.Configuration = config
	if env.Database, err = sqlx.Open(config.Database.Type, config.Database.Url); err != nil {
		return
	}
	if env.Repositories.Images, err = repo.NewImageRepo(env.Database); err != nil {
		return
	}
	if env.Repositories.Albums, err = repo.NewAlbumRepo(env.Database); err != nil {
		return
	}
	if env.Repositories.AlbumImages, err = repo.NewAlbumImageRepo(env.Database); err != nil {
		return
	}
	if env.Storage, err = storage.NewStorage(env.Configuration.Storage); err != nil {
		return
	}
	env.QueueClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
	})
	env.QueueClient.SetDefaultOptions(
		config.Conversion.TaskId,
		asynq.MaxRetry(config.Conversion.MaxRetry),
		asynq.Timeout(time.Duration(config.Conversion.Timeout)),
		asynq.Queue(config.Conversion.Queue),
		asynq.Unique(time.Duration(config.Conversion.UniqueTimeout)),
	)
	return env, err
}

func (env FrontendEnvironment) Destroy() error {
	if err := env.Database.Close(); err != nil {
		return err
	}
	if err := env.QueueClient.Close(); err != nil {
		return err
	}
	return nil
}
