package repo

import (
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/jmoiron/sqlx"
)

type Images struct {
	db              *sqlx.DB
	queryList       *sqlx.NamedStmt
	queryGet        *sqlx.NamedStmt
	stmtCreate      *sqlx.NamedStmt
	stmtUpdate      *sqlx.NamedStmt
	stmtUpdateState *sqlx.NamedStmt
	stmtDelete      *sqlx.NamedStmt
}

const (
	StateCreated    = "created"
	StateQueued     = "queued"
	StateInProgress = "in_progress"
	StateDone       = "done"
	StateError      = "error"
)

func NewImageRepo(db *sqlx.DB) (repo Images, err error) {
	repo.db = db
	repo.queryList, err = db.PrepareNamed(`
			SELECT id,
			       owner,
			       title,
			       description,
			       original_name,
			       created_at,
			       updated_at,
			       state
			FROM images
			WHERE owner = :userId
			ORDER BY created_at DESC
		`)
	if err != nil {
		return
	}
	repo.queryGet, err = db.PrepareNamed(`
			SELECT id,
			       owner,
			       title,
			       description,
			       original_name,
			       created_at,
			       updated_at,
			       state
			FROM images
			WHERE id = :imageId
		`)
	if err != nil {
		return
	}
	repo.stmtCreate, err = db.PrepareNamed(`
			INSERT INTO images (id, owner, title, description, original_name, type, created_at, updated_at, state)
			VALUES (:imageId, :userId, :title, :description, :originalName, :mimeType, NOW(), NOW(), :state)
		`)
	if err != nil {
		return
	}
	repo.stmtUpdate, err = db.PrepareNamed(`
			UPDATE images 
			SET title = :title, 
			    description = :description, 
			    updated_at = NOW()
			WHERE id = :imageId
		`)
	if err != nil {
		return
	}
	repo.stmtUpdateState, err = db.PrepareNamed(`
			UPDATE images
			SET state = :state
			WHERE id = :imageId
		`)
	if err != nil {
		return
	}
	repo.stmtDelete, err = db.PrepareNamed(`
			DELETE FROM images
			WHERE id = :imageId
		`)
	if err != nil {
		return
	}

	return repo, nil
}

func (repo Images) List(user model.User) (images []model.Image, err error) {
	rows, err := repo.queryList.Queryx(map[string]interface{}{
		"userId": user.Id,
	})
	if err != nil {
		return
	}
	for rows.Next() {
		var image model.Image
		err = rows.StructScan(&image)
		if err != nil {
			return
		}
		images = append(images, image)
	}
	return
}

func (repo Images) Get(imageId string) (image model.Image, err error) {
	err = repo.queryGet.Get(&image, map[string]interface{}{
		"imageId": imageId,
	})
	return
}

func (repo Images) Create(new model.Image) (err error) {
	_, err = repo.stmtCreate.Exec(map[string]interface{}{
		"imageId":      new.Id,
		"userId":       new.Owner,
		"title":        new.Title,
		"description":  new.Description,
		"originalName": new.OriginalName,
		"mimeType":     new.MimeType,
		"state":        StateCreated,
	})
	return
}

func (repo Images) Update(changed model.Image) (err error) {
	_, err = repo.stmtUpdate.Exec(map[string]interface{}{
		"imageId":     changed.Id,
		"title":       changed.Title,
		"description": changed.Description,
	})
	return
}

func (repo Images) UpdateState(imageId string, state string) (err error) {
	_, err = repo.stmtUpdateState.Exec(map[string]interface{}{
		"imageId": imageId,
		"state":   state,
	})
	return
}

func (repo Images) Delete(changed model.Image) (err error) {
	_, err = repo.stmtDelete.Exec(map[string]interface{}{
		"imageId": changed.Id,
	})
	return
}
