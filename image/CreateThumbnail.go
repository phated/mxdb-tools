package image

import (
	"mxdb-tools/csv"
	"mxdb-tools/fs"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func CreateThumbnail(card *csv.Card) error {
	largePath := filepath.Join(dirs.Large, card.Filename())
	thumbnailPath := filepath.Join(dirs.Thumbnail, card.Filename())

	if fs.Exists(thumbnailPath) {
		return nil
	}

	img, imgErr := imaging.Open(largePath)
	if imgErr != nil {
		return imgErr
	}

	height := 100
	resizedImg := imaging.Resize(img, 0, height, imaging.Box)

	return imaging.Save(resizedImg, thumbnailPath)
}
