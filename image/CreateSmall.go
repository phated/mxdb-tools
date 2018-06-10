package image

import (
	"mxdb-tools/csv"
	"mxdb-tools/fs"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func CreateSmall(card *csv.Card) error {
	largePath := filepath.Join(dirs.Large, card.Filename())
	smallPath := filepath.Join(dirs.Small, card.Filename())

	if fs.Exists(smallPath) {
		return nil
	}

	img, imgErr := imaging.Open(largePath)
	if imgErr != nil {
		return imgErr
	}

	height := 200
	resizedImg := imaging.Resize(img, 0, height, imaging.Box)

	return imaging.Save(resizedImg, smallPath)
}
