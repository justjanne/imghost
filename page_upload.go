package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

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

func createImage(config *Config, body io.ReadCloser, fileHeader *multipart.FileHeader) (Image, error) {
	id := generateId()
	path := filepath.Join(config.SourceFolder, id)

	err := writeBody(body, path)
	if err != nil {
		return Image{}, err
	}

	mimeType, err := detectMimeType(path)
	if err != nil {
		return Image{}, err
	}

	image := Image{
		Id:           id,
		OriginalName: filepath.Base(fileHeader.Filename),
		CreatedAt:    time.Now(),
		MimeType:     mimeType,
	}
	return image, nil
}

func pageUpload(ctx PageContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			user := parseUser(r)

			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				formatError(w, ErrorData{500, user, r.URL, err}, "json")
				return
			}

			file, header, err := r.FormFile("file")
			if err != nil {
				formatError(w, ErrorData{500, user, r.URL, err}, "json")
				return
			}
			image, err := createImage(ctx.Config, file, header)
			if err != nil {
				formatError(w, ErrorData{500, user, r.URL, err}, "json")
				return
			}

			pubsub := ctx.Redis.Subscribe(ctx.Context, ctx.Config.ResultChannel)
			if _, err = ctx.Database.Exec("INSERT INTO images (id, owner, created_at, updated_at, original_name, type) VALUES ($1, $2, $3, $4, $5, $6)", image.Id, user.Id, image.CreatedAt, image.CreatedAt, image.OriginalName, image.MimeType); err != nil {
				formatError(w, ErrorData{500, user, r.URL, err}, "json")
				return
			}

			data, err := json.Marshal(image)
			if err != nil {
				formatError(w, ErrorData{500, user, r.URL, err}, "json")
				return
			}

			fmt.Printf("Created task %s at %d\n", image.Id, time.Now().Unix())
			ctx.Redis.RPush(ctx.Context, fmt.Sprintf("queue:%s", ctx.Config.ImageQueue), data)
			fmt.Printf("Submitted task %s at %d\n", image.Id, time.Now().Unix())

			waiting := true
			for waiting {
				message, err := pubsub.ReceiveMessage(ctx.Context)
				if err != nil {
					formatError(w, ErrorData{500, user, r.URL, err}, "json")
					return
				}

				result := Result{}
				err = json.Unmarshal([]byte(message.Payload), &result)
				if err != nil {
					formatError(w, ErrorData{500, user, r.URL, err}, "json")
					return
				}

				fmt.Printf("Returned task %s at %d\n", result.Id, time.Now().Unix())

				if result.Id == image.Id {
					waiting = false

					if err = returnJson(w, result); err != nil {
						formatError(w, ErrorData{500, user, r.URL, err}, "json")
						return
					}
				}
			}
			return
		} else {
			user := parseUser(r)
			if err := formatTemplate(w, "upload.html", IndexData{
				user,
			}); err != nil {
				formatError(w, ErrorData{500, user, r.URL, err}, "html")
				return
			}
		}
	})
}
