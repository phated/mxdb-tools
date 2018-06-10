package gql

import (
	"encoding/json"
)

// FetchCards fetches all cards from the API
func FetchCards() ([]*Card, error) {
	type allCards struct {
		AllCards []*Card `json:"allCards"`
	}

	query, err := queries.MustBytes("AllData.graphql")

	if err != nil {
		return nil, err
	}

	respBody, readErr := Request(query, nil)
	if readErr != nil {
		return nil, readErr
	}

	jsonResp := &allCards{}
	if err := json.Unmarshal(respBody, jsonResp); err != nil {
		return nil, err
	}

	return jsonResp.AllCards, nil
}

// FetchTraits fetches all traits from the API
func FetchTraits() ([]*Trait, error) {
	type allTraits struct {
		AllTraits []*Trait `json:"allTraits"`
	}

	query, err := queries.MustBytes("AllTraits.graphql")
	if err != nil {
		return nil, err
	}

	respBody, readErr := Request(query, nil)
	if readErr != nil {
		return nil, readErr
	}

	jsonResp := &allTraits{}
	if err := json.Unmarshal(respBody, jsonResp); err != nil {
		return nil, err
	}

	return jsonResp.AllTraits, nil
}

// FetchStatRanks fetches all stats by type
func FetchStatsByType() (*StatRanks, error) {
	query, err := queries.MustBytes("StatRanks.graphql")
	if err != nil {
		return nil, err
	}

	respBody, readErr := Request(query, nil)
	if readErr != nil {
		return nil, readErr
	}

	jsonResp := &StatRanks{}
	if err := json.Unmarshal(respBody, jsonResp); err != nil {
		return nil, err
	}

	return jsonResp, nil
}
