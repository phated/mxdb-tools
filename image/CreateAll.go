package image

import (
	"mxdb-tools/csv"
)

func CreateAll(card *csv.Card) error {
	if err := CreateOriginal(card); err != nil {
		return err
	}
	if err := CreateLarge(card); err != nil {
		return err
	}
	if err := CreateMedium(card); err != nil {
		return err
	}
	if err := CreateSmall(card); err != nil {
		return err
	}
	if err := CreateThumbnail(card); err != nil {
		return err
	}

	return nil
}
