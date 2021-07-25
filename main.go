package main

import (
	"fmt"
	"git.kuschku.de/justjanne/imghost-frontend/environment"
	"git.kuschku.de/justjanne/imghost-frontend/model"
	"git.kuschku.de/justjanne/imghost-frontend/repo"
	_ "github.com/lib/pq"
)

func main() {
	user := model.User{
		Id:    "ad45284c-be4d-4546-8171-41cf126ac091",
		Name:  "justJanne",
		Email: "janne@kuschku.de",
		Roles: []string{"imghost:user", "imghost:admin"},
	}

	var env environment.Environment
	env, err := environment.InitializeEnvironment()
	if err != nil {
		panic(err)
	}
	defer env.Destroy()

	imageRepo, err := repo.NewImageRepo(env.Database)
	if err != nil {
		panic(err)
	}

	albumRepo, err := repo.NewAlbumRepo(env.Database)
	if err != nil {
		panic(err)
	}

	albumImageRepo, err := repo.NewAlbumImageRepo(env.Database)
	if err != nil {
		panic(err)
	}

	images, err := imageRepo.List(user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("images: %v\n", len(images))

	albums, err := albumRepo.List(user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("albums: %v\n", len(albums))

	albumImages, err := albumImageRepo.List(model.Album{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("albumImages: %v\n", len(albumImages))
}
