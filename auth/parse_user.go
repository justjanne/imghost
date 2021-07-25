package auth

import (
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"net/http"
)

func ParseUser(request *http.Request, env environment.Environment) (user model.User, err error) {
	// TODO: Implement actual user auth
	user = model.User{
		Id:    "ad45284c-be4d-4546-8171-41cf126ac091",
		Name:  "justJanne",
		Email: "janne@kuschku.de",
		Roles: []string{"imghost:user", "imghost:admin"},
	}

	return
}
