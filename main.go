package main

import (
	"flag"
	"log"
	"mxdb-tools/csv"
	"mxdb-tools/gql"
	"mxdb-tools/image"
)

var token string
var dropboxDir string

func init() {
	flag.StringVar(&token, "token", "", "Pass the token for the graphql API")
	flag.StringVar(&dropboxDir, "dropbox", "", "Dropbox directory where large images are copied")
}

func main() {
	flag.Parse()

	if token == "" {
		log.Println("Token required. Use --token")
		return
	} else {
		gql.SetToken(token)
	}

	image.SetDropboxDir(dropboxDir)

	cards, err := csv.Fetch()
	if err != nil {
		log.Println(err)
		return
	}

	for _, card := range cards {
		if err := image.CreateOriginal(card); err != nil {
			log.Println(err)
			continue
		}

		if err := image.CreateLarge(card); err != nil {
			log.Println(err)
			continue
		}
		if err := image.CreateMedium(card); err != nil {
			log.Println(err)
			continue
		}
		if err := image.CreateSmall(card); err != nil {
			log.Println(err)
			continue
		}
		if err := image.CreateThumbnail(card); err != nil {
			log.Println(err)
			continue
		}
	}

	gqlCards, err := gql.FetchCards()
	if err != nil {
		log.Println(err)
		return
	}
	// TODO: Should this be the output of loadGraphQL?
	currentCards := make(map[string]*gql.Card)
	for _, gqlCard := range gqlCards {
		currentCards[gqlCard.UID] = gqlCard
	}

	var createCards []*csv.Card
	for _, card := range cards {
		currentCard := currentCards[card.UID]

		if currentCard == nil {
			createCards = append(createCards, card)
			continue
		}

		if currentCard.Preview.IsEqual(card) == false {
			resp, err := currentCard.Preview.Update(card)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Preview updated: %s", resp)
			}
		}

		if currentCard.Image.IsEqual(card) == false {
			var resp []byte
			var err error
			if currentCard.Image.IsEmpty() {
				resp, err = currentCard.CreateImage(card)
			} else {
				resp, err = currentCard.Image.Update(card)
			}
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Image updated: %s", resp)
			}
		}

		if currentCard.IsEqual(card) == false {
			resp, err := currentCard.Update(card)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Card updated: %s", resp)
			}
		}
	}

	for _, card := range createCards {
		respBody, err := gql.CreateCard(card)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%s", respBody)
	}
}
