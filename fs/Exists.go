package fs

import "os"

func Exists(path string) bool {
	_, statErr := os.Stat(path)
	return os.IsNotExist(statErr) == false
}
