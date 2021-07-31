package environment

import (
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
	"git.kuschku.de/justjanne/imghost-frontend/storage"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
)

type BackendEnvironment struct {
	Configuration configuration.BackendConfiguration
	QueueServer   *asynq.Server
	Database      *sqlx.DB
	Repositories  Repositories
	Storage       storage.Storage
}

func NewBackendEnvironment(config configuration.BackendConfiguration) (env BackendEnvironment, err error) {
	env.Configuration = config
	if env.Database, err = sqlx.Open(config.Database.Type, config.Database.Url); err != nil {
		return
	}
	if env.Repositories.Images, err = repo.NewImageRepo(env.Database); err != nil {
		return
	}
	if env.Repositories.ImageMetadata, err = repo.NewImageMetadataRepo(env.Database); err != nil {
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
	env.QueueServer = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.Redis.Address,
			Password: config.Redis.Password,
		},
		asynq.Config{
			Concurrency:    config.Queue.Concurrency,
			LogLevel:       asynq.LogLevel(config.Queue.LogLevel),
			Queues:         config.Queue.Queues,
			StrictPriority: config.Queue.StrictPriority,
		},
	)
	return env, err
}

func (env BackendEnvironment) Destroy() error {
	if err := env.Database.Close(); err != nil {
		return err
	}
	env.QueueServer.Shutdown()
	return nil
}
