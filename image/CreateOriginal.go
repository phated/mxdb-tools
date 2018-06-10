package image

import (
	"errors"
	"io"
	"log"
	"mxdb-tools/csv"
	"mxdb-tools/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateOriginal(card *csv.Card) error {
	if card.OriginalImageURL == "" {
		return errors.New("Missing Original Image for URL: " + card.UID)
	}

	path := filepath.Join(dirs.Original, card.Filename())

	if fs.Exists(path) {
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
