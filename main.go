package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gocarina/gocsv"
)

const (
	csvUrl     = "https://docs.google.com/spreadsheets/d/1w2TuX7u_wdxFXnUWb_KyRS6o_8vxAEjZV5u5BpkOuI0/export?exportFormat=csv"
	graphqlUrl = "https://api.graph.cool/simple/v1/metaxdb"
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

var dirs ImageDirectories

func init() {
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

func (card Card) Filename() string {
	return card.UID + ".jpg"
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

func fileExists(path string) bool {
	_, statErr := os.Stat(path)
	return os.IsNotExist(statErr) == false
}

func createOriginalImage(card *Card) error {
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

	return nil
}

// 1000px height It will probably do some "sips" stuff too
func createLargeImage(card *Card) error {
	ogPath := filepath.Join(dirs.Original, card.Filename())
	largePath := filepath.Join(dirs.Large, card.Filename())

	if fileExists(largePath) {
		return nil
	}

	// TODO: This is probably better done after we fetch
	sips := exec.Command("sips", "--matchTo", "/System/Library/ColorSync/Profiles/Generic RGB Profile.icc", ogPath)
	if err := sips.Run(); err != nil {
		log.Println("Unable to color correct:", ogPath, "- Proceeding...")
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

	return imaging.Save(resizedImg, largePath)
}

func createMediumImage(card *Card) error {
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

func createSmallImage(card *Card) error {
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

func createThumbnailImage(card *Card) error {
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
	cards, err := loadCSV()
	// cards, err := loadGraphQL()
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

		// log.Println(originalImageFilepath, "does not exist")
		// TODO: Match CSV and Graphql UIDs and deep equal compare?
		// log.Println(card.toCard())
	}
}
