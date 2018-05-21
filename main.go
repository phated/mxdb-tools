package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gocarina/gocsv"
)

const (
	csvUrl     = "https://docs.google.com/spreadsheets/d/1w2TuX7u_wdxFXnUWb_KyRS6o_8vxAEjZV5u5BpkOuI0/export?exportFormat=csv"
	graphqlUrl = "https://api.graph.cool/simple/v1/metaxdb"
)

type Card struct {
	UID               string `csv:"uid"`
	Rarity            string `csv:"rarity"`
	Number            int    `csv:"number"`
	Set               string `csv:"set"`
	Title             string `csv:"title"`
	Subtitle          string `csv:"subtitle"`
	Type              string `csv:"type"`
	Trait             string `csv:"trait"`
	MP                int    `csv:"mp"`
	Symbol            string `csv:"symbol"`
	Effect            string `csv:"effect"`
	Strength          int    `csv:"strength"`     // TODO: This are being defaulted as 0
	Intelligence      int    `csv:"intelligence"` // TODO: This are being defaulted as 0
	Special           int    `csv:"special"`      // TODO: This are being defaulted as 0
	PreviewURL        string `csv:"preview_url"`
	Previewer         string `csv:"previewer"`
	OriginalImageURL  string `csv:"original_image_url"`
	LargeImageURL     string `csv:"large_image_url"`
	SmallImageURL     string `csv:"small_image_url"`
	MediumImageURL    string `csv:"medium_image_url"`
	ThumbnailImageURL string `csv:"thumbnail_image_url"`
}

func loadCSV() ([]*Card, error) {
	resp, err := http.Get(csvUrl)
	if err != nil {
		// log.Println(err)
		return nil, err
	}

	cards := []*Card{}
	if err := gocsv.Unmarshal(resp.Body, &cards); err != nil {
		// log.Println(err)
		return nil, err
	}

	return cards, nil
}

func queryToRequest(queryString []byte) []byte {
	replacer := strings.NewReplacer("\n", "")
	compactQuery := replacer.Replace(string(queryString))
	return []byte(`{"query":"` + compactQuery + `"}`)
}

type GraphqlTrait struct {
	Name string `json:"name"`
}

type GraphqlEffect struct {
	Symbol string `json:"symbol"`
	Text   string `json:"text"`
}

type GraphqlStats struct {
	Type string `json:"type"`
	Rank int    `json:"rank"`
}

type GraphqlImage struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Thumbnail string `json:"thumbnail"`
}

type GraphqlPreview struct {
	Previewer  string `json:"previewer"`
	PreviewURL string `json:"previewUrl"`
	IsActive   bool   `json:"isActive"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type GraphqlCard struct {
	UID       string         `json:"uid"`
	Rarity    string         `json:"rarity"`
	Number    int            `json:"number"`
	Set       string         `json:"set"`
	Title     string         `json:"title"`
	Subtitle  string         `json:"subtitle"`
	Type      string         `json:"type"`
	Trait     GraphqlTrait   `json:"trait"`
	MP        int            `json:"mp"`
	Effect    GraphqlEffect  `json:"effect"`
	Stats     []GraphqlStats `json:"stats"`
	ImageURL  string         `json:"imageUrl"`
	Image     GraphqlImage   `json:"image"`
	Preview   GraphqlPreview `json:"preview"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
}

func (gqlCard *GraphqlCard) toCard() Card {
	card := Card{
		UID:        gqlCard.UID,
		Rarity:     gqlCard.Rarity,
		Number:     gqlCard.Number,
		Set:        gqlCard.Set,
		Title:      gqlCard.Title,
		Subtitle:   gqlCard.Subtitle,
		Type:       gqlCard.Type,
		Trait:      gqlCard.Trait.Name,
		MP:         gqlCard.MP,
		Symbol:     gqlCard.Effect.Symbol,
		Effect:     gqlCard.Effect.Text,
		PreviewURL: gqlCard.Preview.PreviewURL,
		Previewer:  gqlCard.Preview.Previewer,
	}

	for _, stat := range gqlCard.Stats {
		switch stat.Type {
		case "Strength":
			card.Strength = stat.Rank
		case "Intelligence":
			card.Intelligence = stat.Rank
		case "Special":
			card.Special = stat.Rank
		}
	}

	return card
}

func loadGraphQL() ([]*GraphqlCard, error) {
	type AllCards struct {
		AllCards []*GraphqlCard `json:"allCards"`
	}

	type GraphQlResponse struct {
		Cards AllCards `json:"data"`
	}

	queryStr, err := ioutil.ReadFile("queries/AllData.graphql")

	if err != nil {
		return nil, err
	}

	str := queryToRequest(queryStr)
	body := bytes.NewBuffer(str)

	resp, err := http.Post(graphqlUrl, "application/json", body)
	if err != nil {
		return nil, err
	}

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	jsonResp := GraphQlResponse{}
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		return nil, err
	}

	return jsonResp.Cards.AllCards, nil
}

func main() {
	// cards, err := loadCSV()
	cards, err := loadGraphQL()
	if err != nil {
		log.Println(err)
		return
	}

	for _, card := range cards {
		// TODO: Match CSV and Graphql UIDs and deep equal compare?
		log.Println(card.toCard())
	}
}
