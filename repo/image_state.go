package repo

import (
	"github.com/jmoiron/sqlx"
)

const (
	StateCreated    = "created"
	StateQueued     = "queued"
	StateInProgress = "in_progress"
	StateDone       = "done"
	StateError      = "error"
)

type ImageStates struct {
	db         *sqlx.DB
	queryGet   *sqlx.NamedStmt
	stmtUpdate *sqlx.NamedStmt
}

func NewImageStateRepo(db *sqlx.DB) (repo ImageStates, err error) {
	repo.db = db
	repo.queryGet, err = db.PrepareNamed(`
			SELECT state
			FROM images
			WHERE id = :imageId
		`)
	if err != nil {
		return
	}
	repo.stmtUpdate, err = db.PrepareNamed(`
			UPDATE images
			SET state = :state
			WHERE id = :imageId
		`)
	if err != nil {
		return
	}

	return repo, nil
}

func (repo ImageStates) Get(imageId string) (state string, err error) {
	err = repo.queryGet.Get(&state, map[string]interface{}{
		"imageId": imageId,
	})
	return
}

func (repo ImageStates) Update(imageId string, state string) (err error) {
	_, err = repo.stmtUpdate.Exec(map[string]interface{}{
		"imageId": imageId,
		"state":   state,
	})
	return
}
