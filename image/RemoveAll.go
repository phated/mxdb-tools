package image

import (
	"mxdb-tools/csv"
	"os"
	"path/filepath"
)

func RemoveAll(card *csv.Card) error {
	ogPath := filepath.Join(dirs.Original, card.Filename())
	largePath := filepath.Join(dirs.Large, card.Filename())
	mediumPath := filepath.Join(dirs.Medium, card.Filename())
	smallPath := filepath.Join(dirs.Small, card.Filename())
	thumbnailPath := filepath.Join(dirs.Thumbnail, card.Filename())

	if err := os.Remove(ogPath); err != nil {
		return err
	}

	if err := os.Remove(largePath); err != nil {
		return err
	}

	if err := os.Remove(mediumPath); err != nil {
		return err
	}

	if err := os.Remove(smallPath); err != nil {
		return err
	}

	if err := os.Remove(thumbnailPath); err != nil {
		return err
	}

	return nil
}
