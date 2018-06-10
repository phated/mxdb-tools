package gql

import (
	"log"

	"github.com/gobuffalo/packr"
)

var queries packr.Box
var strengthIDs = make(map[int]string)
var intelligenceIDs = make(map[int]string)
var specialIDs = make(map[int]string)
var traitIDs = make(map[string]string)

func init() {
	queries = packr.NewBox("queries/")

	statRanks, err := FetchStatsByType()
	if err != nil {
		log.Fatal(err)
	}

	for _, stat := range statRanks.Strength {
		strengthIDs[stat.Rank] = stat.ID
	}
	for _, stat := range statRanks.Intelligence {
		intelligenceIDs[stat.Rank] = stat.ID
	}
	for _, stat := range statRanks.Special {
		specialIDs[stat.Rank] = stat.ID
	}

	traits, err := FetchTraits()
	if err != nil {
		log.Fatal(err)
	}

	for _, trait := range traits {
		traitIDs[trait.Name] = trait.ID
	}
}
