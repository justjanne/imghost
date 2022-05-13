package main

import (
	"net/http"
	"path"
)

type AlbumDetailData struct {
	User   UserInfo
	Album  Album
	IsMine bool
}

func pageAlbumDetail(ctx PageContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := parseUser(r)

		_, albumId := path.Split(r.URL.Path)
		result, err := ctx.Database.Query(`
			SELECT
				id,
				owner,
				coalesce(title,  ''),
				coalesce(description, ''),
        		coalesce(created_at, to_timestamp(0))
			FROM albums
			WHERE id = $1
			`, albumId)
		if err != nil {
			formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
			return
		}

		var info Album
		if result.Next() {
			var owner string
			err := result.Scan(&info.Id, &owner, &info.Title, &info.Description, &info.CreatedAt)
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}

			result, err := ctx.Database.Query(`
			SELECT
				image,
				title,
				description,
				position
			FROM album_images
			WHERE album = $1
			ORDER BY position ASC
			`, albumId)
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}

			for result.Next() {
				var image AlbumImage
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
