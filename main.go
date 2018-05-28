package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
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
	csvURL     = "https://docs.google.com/spreadsheets/d/1w2TuX7u_wdxFXnUWb_KyRS6o_8vxAEjZV5u5BpkOuI0/export?exportFormat=csv"
	graphqlURL = "https://api.graph.cool/simple/v1/metaxdb"
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

type GraphqlMutations struct {
	CreateCharacterCardWithPreview []byte
	CreateEventCardWithPreview     []byte
	CreateBattleCardWithPreview    []byte
}

func (m *GraphqlMutations) Prepare() error {
	if query, err := ioutil.ReadFile("queries/CreateCharacterCardWithPreview.graphql"); err != nil {
		return err
	} else {
		m.CreateCharacterCardWithPreview = query
	}
	if query, err := ioutil.ReadFile("queries/CreateEventCardWithPreview.graphql"); err != nil {
		return err
	} else {
		m.CreateEventCardWithPreview = query
	}
	if query, err := ioutil.ReadFile("queries/CreateBattleCardWithPreview.graphql"); err != nil {
		return err
	} else {
		m.CreateBattleCardWithPreview = query
	}

	return nil
}

var token string
var dropboxDir string
var dirs ImageDirectories
var mutations GraphqlMutations
var strengthIDs map[int]string
var intelligenceIDs map[int]string
var specialIDs map[int]string
var traitIDs map[string]string

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

	mutations = GraphqlMutations{}
	if err := mutations.Prepare(); err != nil {
		log.Fatal(err)
	}

	strengthIDs = make(map[int]string)
	strengthIDs[1] = "cjersov3s14fu01832cku6ivc"
	strengthIDs[2] = "cjertinv814oj0183jx107b9u"
	strengthIDs[3] = "cjertith714om0183tr2huxiw"
	strengthIDs[4] = "cjertj05014or0183we74xnsb"
	strengthIDs[5] = "cjertj34914ov0183sl7kmqm1"
	strengthIDs[6] = "cjertj5ol14p10183cmzat5y9"
	strengthIDs[7] = "cjertj8a814p60183av1hlcty"
	strengthIDs[8] = "cjertjbqr14pc0183uo7656p7"
	strengthIDs[9] = "cjertjebc14ph0183025ls377"

	intelligenceIDs = make(map[int]string)
	intelligenceIDs[1] = "cjertjj9i14pm01835g18maaz"
	intelligenceIDs[2] = "cjertjmjt14pq01838b9o4h92"
	intelligenceIDs[3] = "cjertjpvq14px0183ff0j6i8o"
	intelligenceIDs[4] = "cjertjspz14q00183gypv8iaa"
	intelligenceIDs[5] = "cjertjy5i14q70183vq4en94t"
	intelligenceIDs[6] = "cjertk2k914qa0183zxqjv2q3"
	intelligenceIDs[7] = "cjertk64814qe0183s8gtzive"
	intelligenceIDs[8] = "cjertkf4i14qj0183iafpggqn"
	intelligenceIDs[9] = "cjertkih114qm01833ztywguw"

	specialIDs = make(map[int]string)
	specialIDs[1] = "cjertkmfs14qr0183zvz5p90r"
	specialIDs[2] = "cjertkqcd14qv01833umx8bs6"
	specialIDs[3] = "cjertkthm14r10183p8ymqv7a"
	specialIDs[4] = "cjertkxl714r401837ddh78bw"
	specialIDs[5] = "cjertl0i914r80183y161rfvw"
	specialIDs[6] = "cjertl3mj14rc018322rxbnb3"
	specialIDs[7] = "cjertl7pv14rg01835d77dxcm"
	specialIDs[8] = "cjertlb0h14rj0183wfrkgp4x"
	specialIDs[9] = "cjertleik14rm01837rde1qzn"

	traitIDs = make(map[string]string)
	traitIDs["Hero"] = "cjes4rlte16z90183tpe0zwp8"
	traitIDs["Villain"] = "cjes4ro4e16zc0183wukhv8yy"
	traitIDs["Blue Lantern"] = "cjes4rtu116zg0183b2u6zvqh"
	traitIDs["Green Lantern"] = "cjes4s2gj16zk0183uon6d2qy"
	traitIDs["Indigo Tribe"] = "cjes4s9rt16zo0183dblaykto"
	traitIDs["Orange Lantern"] = "cjes4sgpd16zs0183ri5sae4f"
	traitIDs["Red Lantern"] = "cjes4smq216zw01834711alj5"
	traitIDs["Sinestro Corps"] = "cjes4ssvn17000183imycm8h7"
	traitIDs["Star Sapphire"] = "cjes4sz0t17030183x780qp5t"
	traitIDs["White Lantern"] = "cjes4t5r617070183n5f60a53"
	traitIDs["Black Lantern"] = "cjes4tna6170b0183yxydoxc8"
	traitIDs["Scout Regiment"] = "cjes4ufa5170g0183ne60w0n5"
	traitIDs["Garrison Regiment"] = "cjes4ulp5170k0183wz5ewix6"
	traitIDs["Military Police Regiment"] = "cjes4ushb170o0183gqvjd1pr"
	traitIDs["Titan"] = "cjes4v11x170t0183yn2ripoe"
	traitIDs["Cadet Corps"] = "cjes4vgkg170x0183kj04xyew"
	traitIDs["Human"] = "cjf8iuklm49ie0183uw7s0iq4"
	traitIDs["Bat Family"] = "cjhpvakkz3z8w0177mo71ce5w"
	traitIDs["Arkham Inmate"] = "cjhqssk8i46fi01770g3v6q2g"
}

type Card struct {
	UID               string `csv:"uid" json:"uid"`
	Rarity            string `csv:"rarity" json:"rarity"`
	Number            int    `csv:"number" json:"number"`
	Set               string `csv:"set" json:"set"`
	Title             string `csv:"title" json:"title"`
	Subtitle          string `csv:"subtitle" json:"subtitle,omitempty"`
	Type              string `csv:"type" json:"-"`
	Trait             string `csv:"trait" json:"-"`
	MP                int    `csv:"mp" json:"mp"`
	Symbol            string `csv:"symbol" json:"symbol"`
	Effect            string `csv:"effect" json:"effect,omitempty"`
	Strength          int    `csv:"strength" json:"-"`     // TODO: This are being defaulted as 0
	Intelligence      int    `csv:"intelligence" json:"-"` // TODO: This are being defaulted as 0
	Special           int    `csv:"special" json:"-"`      // TODO: This are being defaulted as 0
	PreviewURL        string `csv:"preview_url" json:"previewUrl,omitempty"`
	Previewer         string `csv:"previewer" json:"previewer,omitempty"`
	PreviewActive     bool   `csv:"preview_active" json:"-"` // TODO: Should we add this to json?
	OriginalImageURL  string `csv:"original_image_url" json:"originalImage"`
	LargeImageURL     string `csv:"large_image_url" json:"largeImage"`
	MediumImageURL    string `csv:"medium_image_url" json:"mediumImage"`
	SmallImageURL     string `csv:"small_image_url" json:"smallImage"`
	ThumbnailImageURL string `csv:"thumbnail_image_url" json:"thumbnailImage"`
}

func (card *Card) Filename() string {
	return card.UID + ".jpg"
}

func (card *Card) MarshalJSON() ([]byte, error) {
	type Alias Card
	type CardWithStatIDs struct {
		*Alias
		StatIDs []string `json:"statsIds,omitempty"`
		TraitID string   `json:"traitId,omitempty"`
	}

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

	cardWithStatIDs := CardWithStatIDs{
		Alias:   (*Alias)(card),
		StatIDs: statIDs,
		TraitID: traitID,
	}

	return json.Marshal(cardWithStatIDs)
}

func loadCSV() ([]*Card, error) {
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

func graphqlRequest(query []byte, variables interface{}) ([]byte, error) {
	str, err := queryToRequest(query, variables)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", str)
	body := bytes.NewBuffer(str)

	req, err := http.NewRequest("POST", graphqlURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func queryToRequest(queryString []byte, variables interface{}) ([]byte, error) {
	type payload struct {
		Query     string      `json:"query"`
		Variables interface{} `json:"variables,omitempty"`
	}

	replacer := strings.NewReplacer("\n", "")
	compactQuery := replacer.Replace(string(queryString))

	return json.Marshal(payload{
		Query:     compactQuery,
		Variables: variables,
	})
}

type GraphqlTrait struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GraphqlEffect struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Text   string `json:"text"`
}

type GraphqlStats struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Rank int    `json:"rank"`
}

type GraphqlImage struct {
	ID        string `json:"id"`
	Original  string `json:"original"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Thumbnail string `json:"thumbnail"`
}

type GraphqlPreview struct {
	ID         string `json:"id,omitempty"`
	Previewer  string `json:"previewer,omitempty"`
	PreviewURL string `json:"previewUrl,omitempty"`
	IsActive   bool   `json:"isActive,omitempty"`
}

func (preview *GraphqlPreview) IsEqual(card *Card) bool {
	return (preview.Previewer == card.Previewer &&
		preview.PreviewURL == card.PreviewURL &&
		preview.IsActive == card.PreviewActive)
}

func (preview *GraphqlPreview) Update() ([]byte, error) {
	query, err := ioutil.ReadFile("queries/UpdatePreview.graphql")
	if err != nil {
		return nil, err
	}

	return graphqlRequest(query, preview)
}

type GraphqlCard struct {
	ID        string         `json:"id"`
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

func (gqlCard *GraphqlCard) toCard() *Card {
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

	return &card
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

	respBody, readErr := graphqlRequest(queryStr, nil)
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
func createLargeImage(card *Card) error {
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
	flag.Parse()

	if token == "" {
		log.Println("Token required. Use --token")
		return
	}

	cards, err := loadCSV()
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

	gqlCards, err := loadGraphQL()
	// TODO: Should this be the output of loadGraphQL?
	currentCards := make(map[string]*GraphqlCard)
	for _, gqlCard := range gqlCards {
		currentCards[gqlCard.UID] = gqlCard
	}

	var createCards []*Card
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

	// TODO: Non-Preview creation
	for _, card := range createCards {
		var query []byte
		if card.Type == "Character" {
			query = mutations.CreateCharacterCardWithPreview
		}

		if card.Type == "Event" {
			query = mutations.CreateEventCardWithPreview
		}

		if card.Type == "Battle" {
			query = mutations.CreateBattleCardWithPreview
		}

		respBody, err := graphqlRequest(query, card)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%s", respBody)
	}
}
