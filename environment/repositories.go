package environment

import "git.kuschku.de/justjanne/imghost-frontend/repo"

type Repositories struct {
	Images        repo.Images
	ImageMetadata repo.ImageMetadata
	Albums        repo.Albums
	AlbumImages   repo.AlbumImages
}
