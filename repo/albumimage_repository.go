package repo

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"time"
)

type AlbumImageRepository struct {
	db *sql.DB
}

func NewAlbumImageRepository(db *sql.DB) AlbumImageRepository {
	return AlbumImageRepository{
		db: db,
	}
}

func (repo AlbumImageRepository) List(album model.Album) ([]model.AlbumImage, error) {
	var albums []model.AlbumImage

	result, err := repo.db.Query(`
			SELECT
				image,
				coalesce(title,  ''),
				coalesce(description, '')
			FROM album_images
			WHERE album = $1
			ORDER BY position
			`, album.Id)
	if err != nil {
		return albums, err
	}

	for result.Next() {
		var album model.AlbumImage

		if err := result.Scan(
			&album.Id, &album.Title, &album.Description,
		); err != nil {
			return albums, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func (repo AlbumImageRepository) Get(album model.Album, imageId string) (model.AlbumImage, error) {
	var albumImage model.AlbumImage

	result, err := repo.db.Query(`
			SELECT
				image,
				coalesce(title,  ''),
				coalesce(description, '')
			FROM album_images
			WHERE album = $1
			AND image = $2
			ORDER BY position`,
		album.Id,
		imageId)
	if err != nil {
		return albumImage, err
	}

	if result.Next() {
		if err := result.Scan(
			&albumImage.Id, &albumImage.Title, &albumImage.Description,
		); err != nil {
			return albumImage, err
		}
	}

	return albumImage, nil
}

func (repo AlbumImageRepository) Append(album model.Album, image model.AlbumImage) error {
	if _, err := repo.db.Exec(`
		INSERT INTO album_images (album, image, title, description, position) 
		VALUES ($1, $2, $3, $4, (
		    SELECT COUNT(image)
		    FROM album_images
		    WHERE album = $1
		))`,
		album.Id,
		image.Id,
		image.Title,
		image.Description,
	); err != nil {
		return err
	}
	if _, err := repo.db.Exec(`
		UPDATE albums SET updated_at = $2 WHERE id = $1`,
		album.Id,
		time.Now().UTC(),
	); err != nil {
		return err
	}

	return nil
}

func (repo AlbumImageRepository) Update(album model.Album, image model.AlbumImage) error {
	if _, err := repo.db.Exec(`
		UPDATE album_images 
		SET title = $3, description = $4 
		WHERE album = $1 
	    AND image = $2`,
		album.Id,
		image.Id,
		image.Title,
		image.Description,
	); err != nil {
		return err
	}
	if _, err := repo.db.Exec(`
		UPDATE albums SET updated_at = $2 WHERE id = $1`,
		album.Id,
		time.Now().UTC(),
	); err != nil {
		return err
	}

	return nil
}

func (repo AlbumImageRepository) DeleteAll(album model.Album) error {
	if _, err := repo.db.Exec(
		"DELETE FROM album_images WHERE album = $1",
		album.Id,
	); err != nil {
		return err
	}
	if _, err := repo.db.Exec(`
		UPDATE albums SET updated_at = $2 WHERE id = $1`,
		album.Id,
		time.Now().UTC(),
	); err != nil {
		return err
	}

	return nil
}

func (repo AlbumImageRepository) Delete(album model.Album, image model.AlbumImage) error {
	if _, err := repo.db.Exec(
		"DELETE FROM album_images WHERE album = $1 AND image = $2",
		album.Id,
		image.Id,
	); err != nil {
		return err
	}
	if _, err := repo.db.Exec(`
		UPDATE albums SET updated_at = $2 WHERE id = $1`,
		album.Id,
		time.Now().UTC(),
	); err != nil {
		return err
	}

	return nil
}

func (repo AlbumImageRepository) Reorder(album model.Album, images []model.AlbumImage) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	for index, image := range images {
		if _, err := tx.Exec(`
		UPDATE album_images 
		SET position = $3 
		WHERE album = $1 
	    AND image = $2`,
			album.Id,
			image.Id,
			index,
		); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(`
		UPDATE albums SET updated_at = $2 WHERE id = $1`,
		album.Id,
		time.Now().UTC(),
	); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
