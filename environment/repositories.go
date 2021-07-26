package environment

import "git.kuschku.de/justjanne/imghost-frontend/repo"

type Repositories struct {
	Images      repo.Images
	ImageStates repo.ImageStates
	Albums      repo.Albums
	AlbumImages repo.AlbumImages
}
