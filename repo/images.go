package repo

import (
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/jmoiron/sqlx"
)

type Images struct {
	db         *sqlx.DB
	queryList  *sqlx.NamedStmt
	queryGet   *sqlx.NamedStmt
	stmtCreate *sqlx.NamedStmt
	stmtUpdate *sqlx.NamedStmt
	stmtDelete *sqlx.NamedStmt
}

func NewImageRepo(db *sqlx.DB) (repo Images, err error) {
	repo.db = db
	repo.queryList, err = db.PrepareNamed(`
			SELECT images.id,
			       images.owner,
			       images.title,
			       images.description,
			       images.original_name,
			       images.created_at,
			       images.updated_at
			FROM images
			WHERE images.owner = :userId
			ORDER BY images.created_at DESC
		`)
	repo.queryGet, err = db.PrepareNamed(`
			SELECT images.id,
			       images.owner,
			       images.title,
			       images.description,
			       images.original_name,
			       images.created_at,
			       images.updated_at
			FROM images
			WHERE images.id = :imageId
		`)
	repo.stmtCreate, err = db.PrepareNamed(`
			INSERT INTO images (id, owner, title, description, original_name, type, created_at, updated_at)
			VALUES (:imageId, :userId, :title, :description, :originalName, :mimeType, NOW(), NOW())
		`)
	repo.stmtUpdate, err = db.PrepareNamed(`
			UPDATE images 
			SET images.title = :title, 
			    images.description = :description, 
			    images.updated_at = NOW()
			WHERE images.id = :imageId
		`)
	repo.stmtDelete, err = db.PrepareNamed(`
			DELETE FROM images
			WHERE images.id = :imageId
		`)

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
	_, err = repo.stmtUpdate.Exec(map[string]interface{}{
		"imageId":      new.Id,
		"userId":       new.Owner,
		"title":        new.Title,
		"description":  new.Description,
		"originalName": new.OriginalName,
		"mimeType":     new.MimeType,
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

func (repo Images) Delete(changed model.Image) (err error) {
	_, err = repo.stmtDelete.Exec(map[string]interface{}{
		"imageId": changed.Id,
	})
	return
}
