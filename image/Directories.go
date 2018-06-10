package image

import "os"

type Directories struct {
	Base      string
	Original  string
	Large     string
	Medium    string
	Small     string
	Thumbnail string
	Dropbox   string
}

func (dirs Directories) Create() error {
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
