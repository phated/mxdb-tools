package gql

import (
	"mxdb-tools/csv"
)

type Card struct {
	ID        string  `json:"id"`
	UID       string  `json:"uid"`
	Rarity    string  `json:"rarity"`
	Number    int     `json:"number"`
	Set       string  `json:"set"`
	Title     string  `json:"title"`
	Subtitle  string  `json:"subtitle"`
	Type      string  `json:"type"`
	Trait     Trait   `json:"trait"`
	MP        int     `json:"mp"`
	Effect    Effect  `json:"effect"`
	Stats     []Stats `json:"stats"`
	ImageURL  string  `json:"imageUrl"`
	Image     Image   `json:"image"`
	Preview   Preview `json:"preview"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

func (c *Card) CreateImage(card *csv.Card) ([]byte, error) {
	create := Image{
		CardID:    c.ID,
		Original:  card.OriginalImageURL,
		Large:     card.LargeImageURL,
		Medium:    card.MediumImageURL,
		Small:     card.SmallImageURL,
		Thumbnail: card.ThumbnailImageURL,
	}

	query, err := queries.MustBytes("CreateImage.graphql")
	if err != nil {
		return nil, err
	}

	return Request(query, create)
}

// Update updates the top-level properties of a Card, no children
func (c *Card) Update(card *csv.Card) ([]byte, error) {
	updated := Card{
		ID:       c.ID,
		UID:      card.UID,
		Rarity:   card.Rarity,
		Number:   card.Number,
		Set:      card.Set,
		Title:    card.Title,
		Subtitle: card.Subtitle,
		Type:     card.Type,
		MP:       card.MP,
	}

	query, err := queries.MustBytes("UpdateCard.graphql")
	if err != nil {
		return nil, err
	}

	return Request(query, updated)
}

// IsEqual checks if there are differences between the top-level Card properties and a csv.Card
func (c *Card) IsEqual(card *csv.Card) bool {
	return (c.UID == card.UID &&
		c.Rarity == card.Rarity &&
		c.Number == card.Number &&
		c.Set == card.Set &&
		c.Title == card.Title &&
		c.Subtitle == card.Subtitle &&
		c.Type == card.Type &&
		c.MP == card.MP)
}
