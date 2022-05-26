package shared

import (
	"git.kuschku.de/justjanne/imghost/imgconv"
	"net/http"
	"strings"
	"time"
)

type UserInfo struct {
	Id    string
	Name  string
	Email string
	Roles []string
}

func ParseUser(r *http.Request) UserInfo {
	return UserInfo{
		r.Header.Get("X-Auth-Subject"),
		r.Header.Get("X-Auth-Username"),
		r.Header.Get("X-Auth-Email"),
		strings.Split(r.Header.Get("X-Auth-Roles"), ","),
	}
}

func (info UserInfo) HasRole(role string) bool {
	for _, r := range info.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type Result struct {
	Id       string           `json:"id"`
	Success  bool             `json:"success"`
	Errors   []string         `json:"errors"`
	Metadata imgconv.Metadata `json:"metadata"`
}

type Image struct {
	Id           string `json:"id"`
	Title        string
	Description  string
	CreatedAt    time.Time
	OriginalName string
	MimeType     string `json:"mimeType"`
}

type AlbumImage struct {
	Id          string
	Title       string
	Description string
	Position    int
}

type Album struct {
	Id          string
	Title       string
	Description string
	CreatedAt   time.Time
	Images      []AlbumImage
}
