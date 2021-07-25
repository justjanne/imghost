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
			SELECT album_images.album,
			       album_images.image,
			       album_images.title,
			       album_images.description
			FROM album_images
			WHERE album_images.album = :albumId
			ORDER BY album_images.position
		`)
	repo.queryGet, err = db.PrepareNamed(`
			SELECT album_images.album,
			       album_images.image,
			       album_images.title,
			       album_images.description
			FROM album_images
			WHERE album_images.album = :albumId
			AND album_images.image = :imageId
			ORDER BY album_images.position
		`)
	repo.stmtCreate, err = db.PrepareNamed(`
			INSERT INTO album_images (album, image, title, description, position)
			VALUES (:albumId, :imageId, :title, :description, (
			    SELECT COUNT(album_images.image)
			    FROM album_images
			    WHERE album_images.album = :albumId
			))
		`)
	repo.stmtUpdate, err = db.PrepareNamed(`
			UPDATE album_images 
			SET album_images.title = :title, 
			    album_images.description = :description
			WHERE album_images.album = :albumId
		    AND album_images.image = :imageId
		`)
	repo.stmtDelete, err = db.PrepareNamed(`
			DELETE FROM album_images
			WHERE album_images.album = :albumId
			AND album_images.image = :imageID
		`)
	repo.stmtDeleteAll, err = db.PrepareNamed(`
			DELETE FROM album_images
			WHERE album_images.album = :albumId
		`)
	repo.stmtReorder, err = db.PrepareNamed(`
			UPDATE album_images 
			SET album_images.position = :position
			WHERE album_images.album = :albumId
		    AND album_images.image = :imageId
		`)

	return repo, nil
}

func (repo AlbumImages) List(album model.Album) (images []model.AlbumImage, err error) {
	rows, err := repo.queryList.Queryx(map[string]interface{}{
		"albumId": album.Id,
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

func (repo AlbumImages) Get(album model.Album, imageId string) (image model.AlbumImage, err error) {
	err = repo.queryGet.Get(&image, map[string]interface{}{
		"albumId": album.Id,
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

func (repo AlbumImages) Reorder(changed model.AlbumImage, position int) (err error) {
	_, err = repo.stmtDeleteAll.Exec(map[string]interface{}{
		"albumId":  changed.Album,
		"imageId":  changed.Image,
		"position": position,
	})
	return
}
