package repo

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"time"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) ImageRepository {
	return ImageRepository{
		db: db,
	}
}

func (repo ImageRepository) List(user model.User) ([]model.Image, error) {
	var images []model.Image

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
		return images, err
	}

	for result.Next() {
		var image model.Image

		if err := result.Scan(
			&image.Id, &image.Owner, &image.Title, &image.Description,
			&image.CreatedAt, &image.OriginalName, &image.MimeType,
		); err != nil {
			return images, err
		}
		images = append(images, image)
	}

	return images, nil
}

func (repo ImageRepository) Get(imageId string) (model.Image, error) {
	var image model.Image

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
			WHERE id = $1
			`, imageId)
	if err != nil {
		return image, err
	}

	if result.Next() {
		if err := result.Scan(
			&image.Id, &image.Owner, &image.Title, &image.Description,
			&image.CreatedAt, &image.OriginalName, &image.MimeType,
		); err != nil {
			return image, err
		}
	}

	return image, nil
}

func (repo ImageRepository) Create(image model.Image) error {
	if _, err := repo.db.Exec(`
		INSERT INTO images (id, owner, title, description, created_at, updated_at, original_name, type) 
		VALUES ($1, $2, $3, $4, $5)`,
		image.Id,
		image.Owner,
		image.Title,
		image.Description,
		time.Now().UTC(),
		time.Now().UTC(),
		image.OriginalName,
		image.MimeType,
	); err != nil {
		return err
	}

	return nil
}

func (repo ImageRepository) Update(image model.Image) error {
	if _, err := repo.db.Exec(
		"UPDATE images SET title = $1, description = $2, updated_at = $3 WHERE id = $4",
		image.Title,
		image.Description,
		time.Now().UTC(),
		image.Id,
	); err != nil {
		return err
	}

	return nil
}

func (repo ImageRepository) Delete(image model.Image) error {
	if _, err := repo.db.Exec(
		"DELETE FROM images WHERE id = $1",
		image.Id,
	); err != nil {
		return err
	}

	return nil
}
