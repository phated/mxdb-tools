package gql

import "mxdb-tools/csv"

// Preview represents the Preview object on the Graphql server
type Preview struct {
	ID         string `json:"id,omitempty"`
	Previewer  string `json:"previewer,omitempty"`
	PreviewURL string `json:"previewUrl,omitempty"`
	IsActive   bool   `json:"isActive,omitempty"`
}

// Update updates a Preview
func (preview *Preview) Update(card *csv.Card) ([]byte, error) {
	updated := Preview{
		ID:         preview.ID,
		Previewer:  card.Previewer,
		PreviewURL: card.PreviewURL,
		IsActive:   card.PreviewActive,
	}

	query, err := queries.MustBytes("UpdatePreview.graphql")
	if err != nil {
		return nil, err
	}

	return Request(query, updated)
}

// IsEqual checks if there are differences between the Preview and a csv.Card
func (preview *Preview) IsEqual(card *csv.Card) bool {
	return (preview.Previewer == card.Previewer &&
		preview.PreviewURL == card.PreviewURL &&
		preview.IsActive == card.PreviewActive)
}
