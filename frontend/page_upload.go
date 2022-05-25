package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadData struct {
	BaseUrl string
	User    UserInfo
}

func detectMimeType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func generateId() string {
	buffer := make([]byte, 4)
	rand.Read(buffer)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buffer)
}

func writeBody(reader io.ReadCloser, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}
	return out.Close()
}

func createImage(config *shared.Config, body io.ReadCloser, fileHeader *multipart.FileHeader) (shared.Image, error) {
	id := generateId()
	path := filepath.Join(config.SourceFolder, id)

	err := writeBody(body, path)
	if err != nil {
		return shared.Image{}, err
	}

	mimeType, err := detectMimeType(path)
	if err != nil {
		return shared.Image{}, err
	}

	image := shared.Image{
		Id:           id,
		OriginalName: filepath.Base(fileHeader.Filename),
		CreatedAt:    time.Now(),
		MimeType:     mimeType,
	}
	return image, nil
}

func pageUpload(env PageEnvironment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			user := parseUser(r)

			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}

			file, header, err := r.FormFile("file")
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}
			image, err := createImage(env.Config, file, header)
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}

			if _, err = env.Database.Exec("INSERT INTO images (id, owner, created_at, updated_at, original_name, type) VALUES ($1, $2, $3, $4, $5, $6)", image.Id, user.Id, image.CreatedAt, image.CreatedAt, image.OriginalName, image.MimeType); err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}

			fmt.Printf("created task %s at %d\n", image.Id, time.Now().Unix())
			t, err := shared.NewImageResizeTask(image.Id)
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}
			info, err := env.AsynqClient.Enqueue(t, asynq.Retention(env.UploadTimeout))
			fmt.Printf("submitted task %s at %d\n", image.Id, time.Now().Unix())
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}
			info, err = waitOnTask(env, info, env.UploadTimeout)
			fmt.Printf("got result for task %s at %d\n", image.Id, time.Now().Unix())
			if err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}
			var result shared.Result
			if err := json.Unmarshal(info.Result, &result); err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}
			if err = returnJson(w, result); err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "json")
				return
			}
			return
		} else {
			user := parseUser(r)
			if err := formatTemplate(w, "upload.html", UploadData{
				env.Config.BaseUrl,
				user,
			}); err != nil {
				formatError(w, ErrorData{http.StatusInternalServerError, user, r.URL, err}, "html")
				return
			}
		}
	})
}
