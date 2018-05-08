package main

import (
	"log"
	"net/http"

	"github.com/gocarina/gocsv"
)

const (
	url = "https://docs.google.com/spreadsheets/d/1w2TuX7u_wdxFXnUWb_KyRS6o_8vxAEjZV5u5BpkOuI0/export?exportFormat=csv"
)

type Card struct {
	UID           string `csv:"uid"`
	Rarity        string `csv:"rarity"`
	Number        int    `csv:"number"`
	Set           string `csv:"set"`
	Title         string `csv:"title"`
	Subtitle      string `csv:"subtitle"`
	Type          string `csv:"type"`
	Trait         string `csv:"trait"`
	MP            int    `csv:"mp"`
	Symbol        string `csv:"symbol"`
	Effect        string `csv:"effect"`
	Strength      int    `csv:"strength"`
	Intelligence  int    `csv:"intelligence"`
	Special       int    `csv:"special"`
	PreviewURL    string `csv:"previewUrl"`
	Previewer     string `csv:"previewer"`
	LargeImageURL string `csv:"largeImageUrl"`
}

func main() {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	cards := []*Card{}
	if err := gocsv.Unmarshal(resp.Body, &cards); err != nil {
		log.Println(err)
		return
	}

	for _, card := range cards {
		log.Println(card)
	}
}
