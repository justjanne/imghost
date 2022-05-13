package main

import (
	"fmt"
	"git.kuschku.de/justjanne/imghost/shared"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"path"
)

type ImageDetailData struct {
	User   UserInfo
	Image  shared.Image
	IsMine bool
}

func pageImageDetail(ctx PageContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := parseUser(r)
		_, imageId := path.Split(r.URL.Path)

		result, err := ctx.Database.Query(`
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
			formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
			return
		}

		var info shared.Image

		if result.Next() {
			var owner string
			if err := result.Scan(&info.Id, &owner, &info.Title, &info.Description, &info.CreatedAt, &info.OriginalName, &info.MimeType); err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}

			switch r.PostFormValue("action") {
			case "update":
				if _, err := ctx.Database.Exec(
					"UPDATE images SET title = $1, description = $2 WHERE id = $3 AND owner = $4",
					r.PostFormValue("title"),
					r.PostFormValue("description"),
					info.Id,
					user.Id,
				); err != nil {
					formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
					return
				}
				if r.PostFormValue("from_js") == "true" {
					if err := returnJson(w, true); err != nil {
						formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
						return
					}
				} else {
					http.Redirect(w, r, r.URL.Path, http.StatusFound)
				}
				return
			case "delete":
				if _, err := ctx.Database.Exec("DELETE FROM images WHERE id = $1 AND owner = $2", info.Id, user.Id); err != nil {
					formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
					return
				}
				for _, definition := range ctx.Config.Sizes {
					if err := os.Remove(path.Join(ctx.Config.TargetFolder, fmt.Sprintf("%s%s", info.Id, definition.Suffix))); err != nil {
						formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
						return
					}
				}
				http.Redirect(w, r, "/me/images", http.StatusFound)
				return
			}

			if err = formatTemplate(w, "image_detail.html", ImageDetailData{
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
