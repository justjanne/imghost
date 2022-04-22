package main

import (
	"database/sql"
	"net/http"
	"path"
	"strconv"
)

type ImageListData struct {
	User     UserInfo
	Images   []Image
	Previous int64
	Current  int64
	Next     int64
}

const PageSize = 30

func paginateImageListQuery(ctx PageContext, user UserInfo, offset int64, pageSize int) (*sql.Rows, error) {
	if offset == 0 {
		return ctx.Database.Query(`
			SELECT
				id,
				coalesce(title,  ''),
				coalesce(description, ''),
        		coalesce(created_at, to_timestamp(0)),
				coalesce(original_name, ''),
				coalesce(type, '')
			FROM images
			WHERE owner = $1
			ORDER BY created_at DESC
			LIMIT $2
		`, user.Id, pageSize)
	} else {
		return ctx.Database.Query(`
			SELECT
				id,
				coalesce(title,  ''),
				coalesce(description, ''),
        		coalesce(created_at, to_timestamp(0)),
				coalesce(original_name, ''),
				coalesce(type, '')
			FROM images
			WHERE owner = $1
			ORDER BY created_at DESC
			LIMIT $3
			OFFSET $2
		`, user.Id, offset, pageSize)
	}
}

func pageImageList(ctx PageContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := parseUser(r)
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
			formatError(w, ErrorData{500, user, r.URL, err}, "html")
			return
		}

		var images []Image
		for result.Next() {
			var info Image
			if err := result.Scan(&info.Id, &info.Title, &info.Description, &info.CreatedAt, &info.OriginalName, &info.MimeType); err != nil {
				formatError(w, ErrorData{500, user, r.URL, err}, "html")
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
			formatError(w, ErrorData{500, user, r.URL, err}, "html")
			return
		}
	})
}
