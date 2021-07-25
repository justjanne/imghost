package repo

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"time"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) AlbumRepository {
	return AlbumRepository{
		db: db,
	}
}

func (repo AlbumRepository) List(user model.User) ([]model.Album, error) {
	var albums []model.Album

	result, err := repo.db.Query(`
			SELECT
				id,
			    owner,
				coalesce(title,  ''),
				coalesce(description, ''),
        		coalesce(created_at, to_timestamp(0)),
				coalesce(original_name, ''),
				coalesce(type, '')
			FROM images
			WHERE owner = $1
			ORDER BY created_at DESC
			`, user.Id)
	if err != nil {
		return albums, err
	}

	for result.Next() {
		var album model.Album

		if err := result.Scan(
			&album.Id, &album.Owner, &album.Title, &album.Description,
			&album.CreatedAt,
		); err != nil {
			return albums, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func (repo AlbumRepository) Get(albumId string) (model.Album, error) {
	var album model.Album

	result, err := repo.db.Query(`
			SELECT
				id,
				owner,
				coalesce(title,  ''),
				coalesce(description, ''),
        		coalesce(created_at, to_timestamp(0)),
        		coalesce(updated_at, to_timestamp(0))
			FROM albums
			WHERE id = $1`,
			albumId)
	if err != nil {
		return album, err
	}

	if result.Next() {
		if err := result.Scan(
			&album.Id, &album.Owner, &album.Title, &album.Description,
			&album.CreatedAt, &album.UpdatedAt,
		); err != nil {
			return album, err
		}
	}

	return album, nil
}

func (repo AlbumRepository) Create(album model.Album) error {
	if _, err := repo.db.Exec(`
		INSERT INTO albums (id, owner, title, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)`,
		album.Id,
		album.Owner,
		album.Title,
		album.Description,
		time.Now().UTC(),
		time.Now().UTC(),
	); err != nil {
		return err
	}

	return nil
}

func (repo AlbumRepository) Update(album model.Album) error {
	if _, err := repo.db.Exec(
		"UPDATE albums SET title = $1, description = $2, updated_at = $3 WHERE id = $4",
		album.Title,
		album.Description,
		time.Now().UTC(),
		album.Id,
	); err != nil {
		return err
	}

	return nil
}

func (repo AlbumRepository) Delete(album model.Album) error {
	if _, err := repo.db.Exec(
		"DELETE FROM albums WHERE id = $1",
		album.Id,
	); err != nil {
		return err
	}

	return nil
}
