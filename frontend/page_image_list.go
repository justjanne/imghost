package main

import (
	"database/sql"
	"git.kuschku.de/justjanne/imghost/shared"
	"net/http"
	"path"
	"strconv"
)

type ImageListData struct {
	User     shared.UserInfo
	Images   []shared.Image
	Previous int64
	Current  int64
	Next     int64
}

const PageSize = 30

func paginateImageListQuery(env PageEnvironment, user shared.UserInfo, offset int64, pageSize int) (*sql.Rows, error) {
	if offset == 0 {
		return env.Database.Query(`
			SELECT
				id,
				title,
				description,
        		created_at,
				original_name,
				type
			FROM images
			WHERE owner = $1
			ORDER BY created_at DESC
			LIMIT $2
		`, user.Id, pageSize)
	} else {
		return env.Database.Query(`
			SELECT
				id,
				title,
				description,
        		created_at,
				original_name,
				type
			FROM images
			WHERE owner = $1
			ORDER BY created_at DESC
			LIMIT $3
			OFFSET $2
		`, user.Id, offset, pageSize)
	}
}

func pageImageList(ctx PageEnvironment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := shared.ParseUser(r)
		_, page := path.Split(r.URL.Path)
		var pageNumber int64
		pageNumber, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			pageNumber = 1
		}

		result, err := paginateImageListQuery(
			ctx,
			user,
			(pageNumber-1)*PageSize,
			PageSize,
		)
		if err != nil {
			formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
			return
		}

		var images []shared.Image
		for result.Next() {
			var info shared.Image
			if err := result.Scan(&info.Id, &info.Title, &info.Description, &info.CreatedAt, &info.OriginalName, &info.MimeType); err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}
			images = append(images, info)
		}

		if err = formatTemplate(w, "image_list.html", ImageListData{
			user,
			images,
			pageNumber - 1,
			pageNumber,
			pageNumber + 1,
		}); err != nil {
			formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
			return
		}
	})
}
