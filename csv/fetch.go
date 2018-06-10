package csv

import (
	"net/http"

	"github.com/gocarina/gocsv"
)

const (
	csvURL = "https://docs.google.com/spreadsheets/d/1w2TuX7u_wdxFXnUWb_KyRS6o_8vxAEjZV5u5BpkOuI0/export?exportFormat=csv"
)

// Fetch pulls the csv and generates Cards
func Fetch() ([]*Card, error) {
	resp, err := http.Get(csvURL)
	if err != nil {
		return nil, err
	}

	cards := []*Card{}
	if err := gocsv.Unmarshal(resp.Body, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}
