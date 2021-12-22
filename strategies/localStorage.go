package strategies

import (
	"fmt"
	"io/ioutil"
	"os"

	"angorasix.com/media/config"
)

// LocalStrategy ...
type LocalStrategy struct {
	staticDir string
	uploadDir string
}

// CreateLocalStrategy ...
func CreateLocalStrategy(config *config.ServiceConfig) (StorageStrategy, error) {
	staticDir := fmt.Sprintf("./%s", config.StaticDir)
	uploadDir := fmt.Sprintf("%s/%s/", staticDir, config.UploadDir)

	// Creates the static directory if needed
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// path/to/whatever does not exist

		os.Mkdir(staticDir, os.ModePerm)
	}

	// Creates the upload directory if needed
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		// path/to/whatever does not exist

		os.Mkdir(uploadDir, os.ModePerm)
	}

	strategy := &LocalStrategy{
		staticDir,
		uploadDir,
	}

	return strategy, nil
}

// UploadImage uploads image
func (s *LocalStrategy) UploadImage(img *UploadedImageModel) (string, error) {
	// Creates filename.

	filenameToSave := fmt.Sprintf("%s%s_%s", s.uploadDir, createNowString(), *img.Filename)

	// Write the uploaded file to the destination file
	if err := ioutil.WriteFile(filenameToSave, img.Bytes, 0777); err != nil {
		return "", err
	}

	return filenameToSave, nil
}
