package image

func SetDropboxDir(dropboxDir string) {
	if dirs != nil {
		dirs.Dropbox = dropboxDir
	}
}
