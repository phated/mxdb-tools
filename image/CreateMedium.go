package image

import (
	"mxdb-tools/csv"
	"mxdb-tools/fs"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func CreateMedium(card *csv.Card) error {
	largePath := filepath.Join(dirs.Large, card.Filename())
	mediumPath := filepath.Join(dirs.Medium, card.Filename())

	if fs.Exists(mediumPath) {
		return nil
	}

	img, imgErr := imaging.Open(largePath)
	if imgErr != nil {
		return imgErr
	}

	height := 400
	resizedImg := imaging.Resize(img, 0, height, imaging.Box)

	return imaging.Save(resizedImg, mediumPath)
}
