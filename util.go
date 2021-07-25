package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"github.com/go-redis/redis"
	"html/template"
	"net/http"
)



type PageContext struct {
	Config      *Config
	Redis       *redis.Client
	Database    *sql.DB
	Images      http.Handler
	AssetServer http.Handler
}

func parseUser(r *http.Request) model.User {
	return model.User{
		"ad45284c-be4d-4546-8171-41cf126ac091",
		"justJanne",
		"janne@kuschku.de",
		[]string{"imghost:user", "imghost:admin"},
	}

	/*
		return UserInfo{
			r.Header.Get("X-Auth-Subject"),
			r.Header.Get("X-Auth-Username"),
			r.Header.Get("X-Auth-Email"),
			strings.Split(r.Header.Get("X-Auth-Roles"), ","),
		}
	*/
}

func returnJson(w http.ResponseWriter, data interface{}) error {
	marshalled, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(marshalled); err != nil {
		return err
	}

	return nil
}

func returnError(w http.ResponseWriter, code int, message string) error {
	w.WriteHeader(code)
	if _, err := w.Write([]byte(message)); err != nil {
		return err
	}

	return nil
}

func formatTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	pageTemplate, err := template.ParseFiles(
		/*"templates/_base.html",
		"templates/_header.html",
		"templates/_navigation.html",
		"templates/_footer.html",*/
		fmt.Sprintf("templates/%s", templateName),
	)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")
	err = pageTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *PageContext) getImageList(user model.User) ([]model.Image, error) {
	result, err := ctx.Database.Query(`
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
			`, user.Id)
	if err != nil {
		return nil, err
	}

	var images []model.Image
	for result.Next() {
		var info model.Image

		if err := result.Scan(&info.Id, &info.Title, &info.Description, &info.CreatedAt, &info.OriginalName, &info.MimeType); err != nil {
			return nil, err
		}
		images = append(images, info)
	}

	return images, nil
}
