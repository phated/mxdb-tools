package gql

import (
	"mxdb-tools/csv"
)

type Image struct {
	ID        string `json:"id,omitempty"`
	CardID    string `json:"cardId,omitempty"`
	Original  string `json:"original"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Thumbnail string `json:"thumbnail"`
}

// Update updates an Image
func (image *Image) Update(card *csv.Card) ([]byte, error) {
	updated := Image{
		ID:        image.ID,
		Original:  card.OriginalImageURL,
		Large:     card.LargeImageURL,
		Medium:    card.MediumImageURL,
		Small:     card.SmallImageURL,
		Thumbnail: card.ThumbnailImageURL,
	}

	query, err := queries.MustBytes("UpdateImage.graphql")
	if err != nil {
		return nil, err
	}

	return Request(query, updated)
}

// IsEqual checks if there are differences between the Preview and a csv.Card
func (image *Image) IsEqual(card *csv.Card) bool {
	return (image.Original == card.OriginalImageURL &&
		image.Large == card.LargeImageURL &&
		image.Medium == card.MediumImageURL &&
		image.Small == card.SmallImageURL &&
		image.Thumbnail == card.ThumbnailImageURL)
}

func (image *Image) IsEmpty() bool {
	return (image.Original == "" &&
		image.Large == "" &&
		image.Medium == "" &&
		image.Small == "" &&
		image.Thumbnail == "")
}
