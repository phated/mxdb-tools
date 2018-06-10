package image

import (
	"log"
	"os"
	"path/filepath"
)

var dirs *Directories

func init() {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		log.Fatal(cwdErr)
	}

	dirs = &Directories{}
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
