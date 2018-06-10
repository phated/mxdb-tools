package gql

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
