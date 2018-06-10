package csv

type Card struct {
	UID               string `csv:"uid" json:"uid"`
	Rarity            string `csv:"rarity" json:"rarity"`
	Number            int    `csv:"number" json:"number"`
	Set               string `csv:"set" json:"set"`
	Title             string `csv:"title" json:"title"`
	Subtitle          string `csv:"subtitle" json:"subtitle,omitempty"`
	Type              string `csv:"type" json:"-"`
	Trait             string `csv:"trait" json:"-"`
	MP                int    `csv:"mp" json:"mp"`
	Symbol            string `csv:"symbol" json:"symbol"`
	Effect            string `csv:"effect" json:"effect,omitempty"`
	Strength          int    `csv:"strength" json:"-"`     // TODO: This are being defaulted as 0
	Intelligence      int    `csv:"intelligence" json:"-"` // TODO: This are being defaulted as 0
	Special           int    `csv:"special" json:"-"`      // TODO: This are being defaulted as 0
	PreviewURL        string `csv:"preview_url" json:"previewUrl,omitempty"`
	Previewer         string `csv:"previewer" json:"previewer,omitempty"`
	PreviewActive     bool   `csv:"preview_active" json:"-"` // TODO: Should we add this to json?
	OriginalImageURL  string `csv:"original_image_url" json:"originalImage"`
	LargeImageURL     string `csv:"large_image_url" json:"largeImage"`
	MediumImageURL    string `csv:"medium_image_url" json:"mediumImage"`
	SmallImageURL     string `csv:"small_image_url" json:"smallImage"`
	ThumbnailImageURL string `csv:"thumbnail_image_url" json:"thumbnailImage"`
}

// Filename turns the UID into a image filename
func (card *Card) Filename() string {
	return card.UID + ".jpg"
}
