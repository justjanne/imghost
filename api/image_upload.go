package api

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"git.kuschku.de/justjanne/imghost-frontend/auth"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
	"git.kuschku.de/justjanne/imghost-frontend/task"
	"git.kuschku.de/justjanne/imghost-frontend/util"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func generateId() string {
	token := make([]byte, 4)
	fmt.Printf("%v\n", token)
	n, err := rand.Read(token)
	fmt.Printf("%v\n", token)
	if err != nil {
		panic(err)
	}
	if n != 4 {
		panic(errors.New("not enough bytes read"))
	}
	return strings.TrimSuffix(base32.StdEncoding.EncodeToString(token), "=")
}

func determineMimeType(header string, filename string) string {
	mediaType, _, err := mime.ParseMediaType(header)
	if err == nil {
		return mediaType
	}
	mediaType = mime.TypeByExtension(filepath.Ext(filename))
	return mediaType
}

func determineExtension(filename string, mimeType string) (extension string, err error) {
	extension = filepath.Ext(filename)
	if extension != "" {
		return
	}

	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return
	}
	if len(extensions) == 0 {
		err = errors.New("no extensions for type " + mimeType + " found")
		return
	}
	extension = extensions[0]
	return
}

func UploadImage(env environment.FrontendEnvironment) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, err := auth.ParseUser(request, env)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		println("parsed user: " + user.Name)

		err = request.ParseMultipartForm(1 * 1024 * 1024)
		if err != nil {
			println("could not parse multiline form")
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		println("parsed multiline form")

		for key := range request.MultipartForm.File {
			println("found file: " + key)
		}

		var files []model.Image
		for _, file := range request.MultipartForm.File["images[]"] {
			println("processing file")
			var image model.Image
			image.Id = generateId()
			image.Owner = user.Id
			image.OriginalName = file.Filename
			image.MimeType = determineMimeType(file.Header.Get("Content-Type"), file.Filename)
			err = env.Repositories.Images.Create(image)
			if err != nil {
				println("failed creating image: " + image.Id)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			println("created image: " + image.Id)
			data, err := file.Open()
			if err != nil {
				println("failed opening image: " + file.Filename)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Printf("Read image: %d\n", file.Size)
			err = env.Storage.Upload(
				context.Background(),
				env.Configuration.Storage.ConversionBucket,
				image.Id,
				image.MimeType,
				data)
			if err != nil {
				println("failed uploading image: " + file.Filename)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			println("uploaded image: " + env.Configuration.Storage.Endpoint +
				"/" + env.Configuration.Storage.ConversionBucket +
				"/" + image.Id)
			extension, err := determineExtension(image.OriginalName, image.MimeType)
			if err != nil {
				println("failed to determine extension for file")
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			resizeTask, err := task.NewResizeTask(
				image.Id,
				extension,
				env.Configuration)
			if err != nil {
				println("failed creating resize task")
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			println("created task: " + resizeTask.Type())
			_, err = env.QueueClient.Enqueue(resizeTask)
			if err != nil {
				println("failed enqueuing resize task")
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			println("Enqueued task")
			err = env.Repositories.Images.UpdateState(image.Id, repo.StateQueued)
			if err != nil {
				println("failed updating image state")
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			println("Updated state")

			files = append(files, image)
		}
		fmt.Printf("Processed all files: %d\n", len(files))

		util.ReturnJson(writer, files)
	})
}
