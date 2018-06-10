package image

import (
	"log"
	"mxdb-tools/csv"
	"mxdb-tools/fs"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// 1000px height It will probably do some "sips" stuff too
func CreateLarge(card *csv.Card) error {
	ogPath := filepath.Join(dirs.Original, card.Filename())
	largePath := filepath.Join(dirs.Large, card.Filename())

	if fs.Exists(largePath) {
		return nil
	}

	ogImg, ogImgErr := imaging.Open(ogPath)
	if ogImgErr != nil {
		return ogImgErr
	}

	border := 30
	height := 980 + (border * 2)
	width := 680 + (border * 2)
	croppedImage := imaging.CropCenter(ogImg, width, height)
	resizedImg := imaging.Resize(croppedImage, 0, 1000, imaging.Box)

	if dirs.Dropbox != "" && card.PreviewActive == true {
		dropboxPath := filepath.Join(dirs.Dropbox, card.Filename())
		if err := imaging.Save(resizedImg, dropboxPath); err != nil {
			log.Println("Failed to write card to", dropboxPath)
		}
	}

	return imaging.Save(resizedImg, largePath)
}
