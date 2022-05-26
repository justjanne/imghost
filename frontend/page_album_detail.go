package main

import (
	"git.kuschku.de/justjanne/imghost/shared"
	"net/http"
	"path"
)

type AlbumDetailData struct {
	User   shared.UserInfo
	Album  shared.Album
	IsMine bool
}

func pageAlbumDetail(env PageEnvironment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := shared.ParseUser(r)

		_, albumId := path.Split(r.URL.Path)
		result, err := env.Database.Query(`
			SELECT
				id,
				owner,
				title,
				description,
        		created_at
			FROM albums
			WHERE id = $1
			`, albumId)
		if err != nil {
			formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
			return
		}

		var info shared.Album
		if result.Next() {
			var owner string
			err := result.Scan(&info.Id, &owner, &info.Title, &info.Description, &info.CreatedAt)
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}

			result, err := env.Database.Query(`
			SELECT
				image,
				title,
				description,
				position
			FROM album_images
			WHERE album = $1
			ORDER BY position
			`, albumId)
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}

			for result.Next() {
				var image shared.AlbumImage
				err := result.Scan(&image.Id, &owner, &image.Title, &image.Description, &image.Position)
				if err != nil {
					formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
					return
				}

				info.Images = append(info.Images, image)
			}

			if err = formatTemplate(w, "album_detail.html", AlbumDetailData{
				user,
				info,
				owner == user.Id,
			}); err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}

			return
		}

		if err := returnError(w, http.StatusNotFound, "Image Not Found"); err != nil {
			formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
			return
		}
	})
}
