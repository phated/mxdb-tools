package gql

import (
	"fmt"
	"mxdb-tools/csv"
)

// TODO: Non-Preview creation

func CreateCard(card *csv.Card) ([]byte, error) {
	if card.Type == "Character" {
		return CreateCharacterCard(card)
	}

	if card.Type == "Event" {
		return CreateEventCard(card)
	}

	if card.Type == "Battle" {
		return CreateBattleCard(card)
	}

	return nil, fmt.Errorf("Invalid card type: %s", card.Type)
}

func CreateCharacterCard(card *csv.Card) ([]byte, error) {
	query, err := queries.MustBytes("CreateCharacterCardWithPreview.graphql")
	if err != nil {
		return nil, err
	}

	return Request(query, prepareCard(card))
}

func CreateEventCard(card *csv.Card) ([]byte, error) {
	query, err := queries.MustBytes("CreateEventCardWithPreview.graphql")
	if err != nil {
		return nil, err
	}

	return Request(query, prepareCard(card))
}

func CreateBattleCard(card *csv.Card) ([]byte, error) {
	query, err := queries.MustBytes("CreateBattleCardWithPreview.graphql")
	if err != nil {
		return nil, err
	}

	return Request(query, prepareCard(card))
}

/* Create utils */

type preparedCard struct {
	*csv.Card
	StatIDs []string `json:"statsIds,omitempty"`
	TraitID string   `json:"traitId,omitempty"`
}

func prepareCard(card *csv.Card) preparedCard {
	var statIDs []string
	if id := strengthIDs[card.Strength]; id != "" {
		statIDs = append(statIDs, id)
	}
	if id := intelligenceIDs[card.Intelligence]; id != "" {
		statIDs = append(statIDs, id)
	}
	if id := specialIDs[card.Special]; id != "" {
		statIDs = append(statIDs, id)
	}

	var traitID string
	if id := traitIDs[card.Trait]; id != "" {
		traitID = id
	}

	return preparedCard{
		Card:    card,
		StatIDs: statIDs,
		TraitID: traitID,
	}
}