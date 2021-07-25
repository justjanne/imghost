package repo

import (
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/jmoiron/sqlx"
)

type Albums struct {
	db         *sqlx.DB
	queryList  *sqlx.NamedStmt
	queryGet   *sqlx.NamedStmt
	stmtUpdate *sqlx.NamedStmt
	stmtCreate *sqlx.NamedStmt
	stmtDelete *sqlx.NamedStmt
}

func NewAlbumRepo(db *sqlx.DB) (repo Albums, err error) {
	repo.db = db
	repo.queryList, err = db.PrepareNamed(`
			SELECT albums.id,
			       albums.owner,
			       albums.title,
			       albums.description,
			       albums.created_at,
			       albums.updated_at
			FROM albums
			WHERE albums.owner = :userId
			ORDER BY albums.created_at DESC
		`)
	repo.queryGet, err = db.PrepareNamed(`
			SELECT albums.id,
			       albums.owner,
			       albums.title,
			       albums.description,
			       albums.created_at,
			       albums.updated_at
			FROM albums
			WHERE albums.id = :albumId
		`)
	repo.stmtCreate, err = db.PrepareNamed(`
			INSERT INTO albums (id, owner, title, description, created_at, updated_at)
			VALUES (:albumId, :userId, :title, :description, NOW(), NOW())
		`)
	repo.stmtUpdate, err = db.PrepareNamed(`
			UPDATE albums 
			SET albums.title = :title, 
			    albums.description = :description, 
			    albums.updated_at = NOW()
			WHERE albums.id = :albumId
		`)
	repo.stmtDelete, err = db.PrepareNamed(`
			DELETE FROM albums
			WHERE albums.id = :albums
		`)

	return repo, nil
}

func (repo Albums) List(user model.User) (albums []model.Album, err error) {
	rows, err := repo.queryList.Queryx(map[string]interface{}{
		"userId": user.Id,
	})
	if err != nil {
		return
	}
	for rows.Next() {
		var album model.Album
		err = rows.StructScan(&album)
		if err != nil {
			return
		}
		albums = append(albums, album)
	}
	return
}

func (repo Albums) Get(albumId string) (album model.Album, err error) {
	err = repo.queryGet.Get(&album, map[string]interface{}{
		"albumId": albumId,
	})
	return
}

func (repo Albums) Create(changed model.Album) (err error) {
	_, err = repo.stmtCreate.Exec(map[string]interface{}{
		"albumId":     changed.Id,
		"userId":      changed.Owner,
		"title":       changed.Title,
		"description": changed.Description,
	})
	return
}

func (repo Albums) Update(changed model.Album) (err error) {
	_, err = repo.stmtUpdate.Exec(map[string]interface{}{
		"albumId":     changed.Id,
		"title":       changed.Title,
		"description": changed.Description,
	})
	return
}

func (repo Albums) Delete(changed model.Album) (err error) {
	_, err = repo.stmtDelete.Exec(map[string]interface{}{
		"albumId": changed.Id,
	})
	return
}
