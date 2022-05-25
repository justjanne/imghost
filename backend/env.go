package main

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost/shared"
)

type ProcessingEnvironment struct {
	Config   *shared.Config
	Database *sql.DB
}
