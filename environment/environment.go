package environment

import (
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"time"
)

type Environment struct {
	Configuration configuration.Configuration
	QueueClient   *asynq.Client
	QueueServer   *asynq.Server
	Database      *sqlx.DB
	Repositories  Repositories
}

func newCommonEnvironment(config configuration.Configuration) (env Environment, err error) {
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
	return
}

func NewClientEnvironment(config configuration.Configuration) (Environment, error) {
	env, err := newCommonEnvironment(config)
	if err != nil {
		return env, err
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

func NewServerEnvironment(config configuration.Configuration) (Environment, error) {
	env, err := newCommonEnvironment(config)
	if err != nil {
		return env, err
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

func (env Environment) Destroy() error {
	if err := env.Database.Close(); err != nil {
		return err
	}
	if env.QueueClient != nil {
		if err := env.QueueClient.Close(); err != nil {
			return err
		}
	}
	if env.QueueServer != nil {
		env.QueueServer.Shutdown()
	}
	return nil
}
