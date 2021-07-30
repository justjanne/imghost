package repo

import (
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/jmoiron/sqlx"
)

type AlbumImages struct {
	db            *sqlx.DB
	queryList     *sqlx.NamedStmt
	queryGet      *sqlx.NamedStmt
	stmtCreate    *sqlx.NamedStmt
	stmtUpdate    *sqlx.NamedStmt
	stmtDelete    *sqlx.NamedStmt
	stmtDeleteAll *sqlx.NamedStmt
	stmtReorder   *sqlx.NamedStmt
}

func NewAlbumImageRepo(db *sqlx.DB) (repo AlbumImages, err error) {
	repo.db = db
	repo.queryList, err = db.PrepareNamed(`
			SELECT album,
			       image,
			       title,
			       description
			FROM album_images
			WHERE album = :albumId
			ORDER BY position
		`)
	if err != nil {
		return
	}
	repo.queryGet, err = db.PrepareNamed(`
			SELECT album,
			       image,
			       title,
			       description
			FROM album_images
			WHERE album = :albumId
			AND image = :imageId
			ORDER BY position
		`)
	if err != nil {
		return
	}
	repo.stmtCreate, err = db.PrepareNamed(`
			INSERT INTO album_images (album, image, title, description, position)
			VALUES (:albumId, :imageId, :title, :description, (
			    SELECT COUNT(image)
			    FROM album_images
			    WHERE album = :albumId
			))
		`)
	if err != nil {
		return
	}
	repo.stmtUpdate, err = db.PrepareNamed(`
			UPDATE album_images 
			SET title = :title, 
			    description = :description
			WHERE album = :albumId
		    AND image = :imageId
		`)
	if err != nil {
		return
	}
	repo.stmtDelete, err = db.PrepareNamed(`
			DELETE FROM album_images
			WHERE album = :albumId
			AND image = :imageID
		`)
	if err != nil {
		return
	}
	repo.stmtDeleteAll, err = db.PrepareNamed(`
			DELETE FROM album_images
			WHERE album = :albumId
		`)
	if err != nil {
		return
	}
	repo.stmtReorder, err = db.PrepareNamed(`
			UPDATE album_images 
			SET position = :position
			WHERE album = :albumId
		    AND image = :imageId
		`)
	if err != nil {
		return
	}

	return repo, nil
}

func (repo AlbumImages) List(albumId string) (images []model.AlbumImage, err error) {
	rows, err := repo.queryList.Queryx(map[string]interface{}{
		"albumId": albumId,
	})
	if err != nil {
		return
	}
	for rows.Next() {
		var image model.AlbumImage
		err = rows.StructScan(&image)
		if err != nil {
			return
		}
		images = append(images, image)
	}
	return
}

func (repo AlbumImages) Get(albumId string, imageId string) (image model.AlbumImage, err error) {
	err = repo.queryGet.Get(&image, map[string]interface{}{
		"albumId": albumId,
		"imageId": imageId,
	})
	return
}

func (repo AlbumImages) Create(new model.AlbumImage) (err error) {
	_, err = repo.stmtCreate.Exec(map[string]interface{}{
		"albumId":     new.Album,
		"imageId":     new.Image,
		"title":       new.Title,
		"description": new.Description,
	})
	return
}

func (repo AlbumImages) Update(changed model.AlbumImage) (err error) {
	_, err = repo.stmtUpdate.Exec(map[string]interface{}{
		"albumId":     changed.Album,
		"imageId":     changed.Image,
		"title":       changed.Title,
		"description": changed.Description,
	})
	return
}

func (repo AlbumImages) Reorder(changed model.AlbumImage, position int) (err error) {
	_, err = repo.stmtReorder.Exec(map[string]interface{}{
		"albumId":  changed.Album,
		"imageId":  changed.Image,
		"position": position,
	})
	return
}

func (repo AlbumImages) Delete(changed model.AlbumImage) (err error) {
	_, err = repo.stmtDelete.Exec(map[string]interface{}{
		"albumId": changed.Album,
		"imageId": changed.Image,
	})
	return
}

func (repo AlbumImages) DeleteAll(changed model.AlbumImage) (err error) {
	_, err = repo.stmtDeleteAll.Exec(map[string]interface{}{
		"albumId": changed.Album,
	})
	return
}
