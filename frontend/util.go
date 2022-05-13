package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"git.kuschku.de/justjanne/imghost/shared"
	"github.com/hibiken/asynq"
	"html/template"
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

func (info UserInfo) HasRole(role string) bool {
	for _, r := range info.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type PageContext struct {
	Context        context.Context
	Config         *shared.Config
	AsynqClient    *asynq.Client
	AsynqInspector *asynq.Inspector
	UploadTimeout  time.Duration
	Database       *sql.DB
	Images         http.Handler
	AssetServer    http.Handler
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

func parseUser(r *http.Request) UserInfo {
	return UserInfo{
		r.Header.Get("X-Auth-Subject"),
		r.Header.Get("X-Auth-Username"),
		r.Header.Get("X-Auth-Email"),
		strings.Split(r.Header.Get("X-Auth-Roles"), ","),
	}
}

func returnJson(w http.ResponseWriter, data interface{}) error {
	marshalled, err := json.Marshal(data)
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
		"templates/_base.html",
		"templates/_header.html",
		"templates/_navigation.html",
		"templates/_footer.html",
		fmt.Sprintf("templates/%s", templateName),
	)
	if err != nil {
		return err
	}

	err = pageTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func waitOnTask(ctx PageContext, info *asynq.TaskInfo, timeout time.Duration) (*asynq.TaskInfo, error) {
	total := time.Duration(0)
	info, err := ctx.AsynqInspector.GetTaskInfo(info.Queue, info.ID)
	if err != nil {
		return nil, err
	}
	// Wait for it being scheduled
	for total < timeout && info.State == asynq.TaskStateScheduled {
		duration := info.NextProcessAt.Sub(time.Now())
		total += duration
		if total < timeout {
			time.Sleep(duration)
		}
		info, err = ctx.AsynqInspector.GetTaskInfo(info.Queue, info.ID)
		if err != nil {
			return nil, err
		}
	}
	// Wait for it being completed
	for total < timeout && info.State != asynq.TaskStateArchived && info.State != asynq.TaskStateCompleted {
		duration := 1 * time.Second
		total += duration
		if total < timeout {
			time.Sleep(duration)
		}
		info, err = ctx.AsynqInspector.GetTaskInfo(info.Queue, info.ID)
		if err != nil {
			return nil, err
		}
	}
	if info.State == asynq.TaskStateCompleted {
		return info, nil
	} else if info.State == asynq.TaskStateArchived {
		return info, fmt.Errorf("error executing task: %s (%s), has status %s", info.Type, info.ID, info.State)
	} else {
		return info, fmt.Errorf("task timed out: %s (%s), has status %s", info.Type, info.ID, info.State)
	}
}
