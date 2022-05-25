package main

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	"net/http"
	"time"
)

type PageEnvironment struct {
	Config         *shared.Config
	AsynqClient    *asynq.Client
	AsynqInspector *asynq.Inspector
	UploadTimeout  time.Duration
	Database       *sql.DB
	Images         http.Handler
	AssetServer    http.Handler
}
