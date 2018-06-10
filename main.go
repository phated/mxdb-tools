package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"mxdb-tools/csv"
	"mxdb-tools/gql"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type ImageDirectories struct {
	Base      string
	Original  string
	Large     string
	Medium    string
	Small     string
	Thumbnail string
}

func (dirs ImageDirectories) Create() error {
	if err := os.MkdirAll(dirs.Original, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(dirs.Large, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(dirs.Medium, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(dirs.Small, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(dirs.Thumbnail, 0700); err != nil {
		return err
	}

	return nil
}

var token string
var dropboxDir string
var dirs ImageDirectories

func init() {
	flag.StringVar(&token, "token", "", "Pass the token for the graphql API")
	flag.StringVar(&dropboxDir, "dropbox", "", "Dropbox directory where large images are copied")

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		log.Fatal(cwdErr)
	}

	dirs = ImageDirectories{}
	dirs.Base = filepath.Join(cwd, "images/")
	dirs.Original = filepath.Join(dirs.Base, "original/")
	dirs.Large = filepath.Join(dirs.Base, "large/")
	dirs.Medium = filepath.Join(dirs.Base, "medium/")
	dirs.Small = filepath.Join(dirs.Base, "small/")
	dirs.Thumbnail = filepath.Join(dirs.Base, "thumbnail/")

	if err := dirs.Create(); err != nil {
		log.Fatal(err)
	}
}

func fileExists(path string) bool {
	_, statErr := os.Stat(path)
	return os.IsNotExist(statErr) == false
}

func createOriginalImage(card *csv.Card) error {
	if card.OriginalImageURL == "" {
		return errors.New("Missing Original Image for URL: " + card.UID)
	}

	path := filepath.Join(dirs.Original, card.Filename())

	if fileExists(path) {
		return nil
	}

	log.Println("Downloading:", card.OriginalImageURL)

	resp, respErr := http.Get(card.OriginalImageURL)
	if respErr != nil {
		return respErr
	}

	defer resp.Body.Close()

	imgFile, imgErr := os.Create(path)
	if imgErr != nil {
		return imgErr
	}

	defer imgFile.Close()

	_, copyErr := io.Copy(imgFile, resp.Body)
	if copyErr != nil {
		return copyErr
	}

	// TODO: This should probably check color profile instead of Card's Set
	// TODO: Would be nice to make this cross platform
	if card.Set == "JL" {
		sips := exec.Command("sips", "--matchTo", "/System/Library/ColorSync/Profiles/Generic RGB Profile.icc", path)
		if err := sips.Run(); err != nil {
			log.Println("Unable to color correct:", path, "- Proceeding...")
		}
	}

	return nil
}

// 1000px height It will probably do some "sips" stuff too
func createLargeImage(card *csv.Card) error {
	ogPath := filepath.Join(dirs.Original, card.Filename())
	largePath := filepath.Join(dirs.Large, card.Filename())

	if fileExists(largePath) {
		return nil
	}

	ogImg, ogImgErr := imaging.Open(ogPath)
	if ogImgErr != nil {
		return ogImgErr
	}

	border := 30
	height := 980 + (border * 2)
	width := 680 + (border * 2)
	croppedImage := imaging.CropCenter(ogImg, width, height)
	resizedImg := imaging.Resize(croppedImage, 0, 1000, imaging.Box)

	if dropboxDir != "" && card.PreviewActive == true {
		dropboxPath := filepath.Join(dropboxDir, card.Filename())
		if err := imaging.Save(resizedImg, dropboxPath); err != nil {
			log.Println("Failed to write card to", dropboxPath)
		}
	}

	return imaging.Save(resizedImg, largePath)
}

func createMediumImage(card *csv.Card) error {
	largePath := filepath.Join(dirs.Large, card.Filename())
	mediumPath := filepath.Join(dirs.Medium, card.Filename())

	if fileExists(mediumPath) {
		return nil
	}

	img, imgErr := imaging.Open(largePath)
	if imgErr != nil {
		return imgErr
	}

	height := 400
	resizedImg := imaging.Resize(img, 0, height, imaging.Box)

	return imaging.Save(resizedImg, mediumPath)
}

func createSmallImage(card *csv.Card) error {
	largePath := filepath.Join(dirs.Large, card.Filename())
	smallPath := filepath.Join(dirs.Small, card.Filename())

	if fileExists(smallPath) {
		return nil
	}

	img, imgErr := imaging.Open(largePath)
	if imgErr != nil {
		return imgErr
	}

	height := 200
	resizedImg := imaging.Resize(img, 0, height, imaging.Box)

	return imaging.Save(resizedImg, smallPath)
}

func createThumbnailImage(card *csv.Card) error {
	largePath := filepath.Join(dirs.Large, card.Filename())
	thumbnailPath := filepath.Join(dirs.Thumbnail, card.Filename())

	if fileExists(thumbnailPath) {
		return nil
	}

	img, imgErr := imaging.Open(largePath)
	if imgErr != nil {
		return imgErr
	}

	height := 100
	resizedImg := imaging.Resize(img, 0, height, imaging.Box)

	return imaging.Save(resizedImg, thumbnailPath)
}

func main() {
	flag.Parse()

	if token == "" {
		log.Println("Token required. Use --token")
		return
	} else {
		gql.SetToken(token)
	}

	cards, err := csv.Fetch()
	if err != nil {
		log.Println(err)
		return
	}

	for _, card := range cards {
		if err := createOriginalImage(card); err != nil {
			log.Println(err)
			continue
		}

		if err := createLargeImage(card); err != nil {
			log.Println(err)
			continue
		}
		if err := createMediumImage(card); err != nil {
			log.Println(err)
			continue
		}
		if err := createSmallImage(card); err != nil {
			log.Println(err)
			continue
		}
		if err := createThumbnailImage(card); err != nil {
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
			currentCard.Preview.Previewer = card.Previewer
			currentCard.Preview.PreviewURL = card.PreviewURL
			currentCard.Preview.IsActive = card.PreviewActive
			resp, err := currentCard.Preview.Update()
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Preivew updated: %s", resp)
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
