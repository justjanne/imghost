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

type UploadData struct {
	User    UserInfo
	Results []Result
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
				if err = returnJson(w, []Result{{
					Success: false,
					Errors:  []string{err.Error()},
				}}); err != nil {
					panic(err)
				}
			}

			file, header, err := r.FormFile("file")
			if err != nil {
				if err = returnJson(w, []Result{{
					Success: false,
					Errors:  []string{err.Error()},
				}}); err != nil {
					panic(err)
				}
				return
			}
			image, err := createImage(ctx.Config, file, header)
			if err != nil {
				if err = returnJson(w, []Result{{
					Success: false,
					Errors:  []string{err.Error()},
				}}); err != nil {
					panic(err)
				}
				return
			}

			pubsub := ctx.Redis.Subscribe(ctx.Config.ResultChannel)
			_, err = ctx.Database.Exec("INSERT INTO images (id, owner, created_at, updated_at, original_name, type) VALUES ($1, $2, $3, $4, $5, $6)", image.Id, user.Id, image.CreatedAt, image.CreatedAt, image.OriginalName, image.MimeType)
			if err != nil {
				panic(err)
			}

			data, err := json.Marshal(image)
			if err != nil {
				if err = returnJson(w, []Result{{
					Success: false,
					Errors:  []string{err.Error()},
				}}); err != nil {
					panic(err)
				}
				return
			}

			fmt.Printf("Created task %s at %d\n", image.Id, time.Now().Unix())
			ctx.Redis.RPush(fmt.Sprintf("queue:%s", ctx.Config.ImageQueue), data)
			fmt.Printf("Submitted task %s at %d\n", image.Id, time.Now().Unix())

			waiting := true
			for waiting {
				message, err := pubsub.ReceiveMessage()
				if err != nil {
					if err = returnJson(w, []Result{{
						Success: false,
						Errors:  []string{err.Error()},
					}}); err != nil {
						panic(err)
					}
					return
				}

				result := Result{}
				err = json.Unmarshal([]byte(message.Payload), &result)
				if err != nil {
					if err = returnJson(w, []Result{{
						Success: false,
						Errors:  []string{err.Error()},
					}}); err != nil {
						panic(err)
					}
					return
				}

				fmt.Printf("Returned task %s at %d\n", result.Id, time.Now().Unix())

				if result.Id == image.Id {
					waiting = false

					if err = returnJson(w, result); err != nil {
						panic(err)
					}
				}
			}
			return
		} else {
			user := parseUser(r)
			if err := formatTemplate(w, "upload.html", UploadData{
				user,
				[]Result{},
			}); err != nil {
				panic(err)
			}
		}
	})
}
