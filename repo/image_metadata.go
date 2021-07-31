package repo

import (
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/jmoiron/sqlx"
)

type ImageMetadata struct {
	db         *sqlx.DB
	queryList  *sqlx.NamedStmt
	stmtCreate *sqlx.NamedStmt
	stmtDelete *sqlx.NamedStmt
}

func NewImageMetadataRepo(db *sqlx.DB) (repo ImageMetadata, err error) {
	repo.db = db
	repo.queryList, err = db.PrepareNamed(`
			SELECT name, content
			FROM image_metadata
			WHERE image = :imageId
		`)
	if err != nil {
		return
	}
	repo.stmtCreate, err = db.PrepareNamed(`
			INSERT INTO image_metadata (image, name, content)
			VALUES (:imageId, :name, :content)
		`)
	if err != nil {
		return
	}
	repo.stmtDelete, err = db.PrepareNamed(`
			DELETE FROM image_metadata
			WHERE image = :imageId
		`)
	if err != nil {
		return
	}

	return repo, nil
}

func (repo ImageMetadata) List(image model.Image) (map[string]string, error) {
	rows, err := repo.queryList.Queryx(map[string]interface{}{
		"imageId": image.Id,
	})
	if err != nil {
		return nil, err
	}
	metadata := make(map[string]string)
	for rows.Next() {
		var key string
		var value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		metadata[key] = value
	}
	return metadata, nil
}

func (repo ImageMetadata) Update(imageId string, metadata map[string]string) (err error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return
	}
	println("Deleting metadata for " + imageId)
	_, err = tx.NamedStmt(repo.stmtDelete).Exec(map[string]interface{}{
		"imageId": imageId,
	})
	if err != nil {
		return
	}
	for key, value := range metadata {
		println("Adding metadata for " + imageId + " with " + key + "=" + value)
		_, err = tx.NamedStmt(repo.stmtCreate).Exec(map[string]interface{}{
			"imageId": imageId,
			"name":    key,
			"content": value,
		})
		if err != nil {
			return
		}
	}
	err = tx.Commit()
	return
}

func (repo ImageMetadata) Delete(changed model.Image) (err error) {
	_, err = repo.stmtDelete.Exec(map[string]interface{}{
		"imageId": changed.Id,
	})
	return
}
