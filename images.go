package media

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/nfnt/resize"

	"angorasix.com/media/config"
	images "angorasix.com/media/gen/images"
	strategies "angorasix.com/media/strategies"
)

// images service example implementation.
// The example methods log the requests and return zero values.
type imagessrvc struct {
	logger   *log.Logger
	strategy strategies.StorageStrategy
	config   *config.ServiceConfig
}

var mapCode = map[int]png.CompressionLevel{
	0:  png.DefaultCompression,
	-1: png.NoCompression,
	-2: png.BestSpeed,
	-3: png.BestCompression,
}

// NewImages returns the images service implementation.
func NewImages(logger *log.Logger, strategy strategies.StorageStrategy, config *config.ServiceConfig) images.Service {
	return &imagessrvc{logger, strategy, config}
}

// Uploads a new image
func (s *imagessrvc) Upload(ctx context.Context, p *images.ImageUploadPayload) (res *images.ImagesListMedia, err error) {
	res = &images.ImagesListMedia{}
	imagesCollection := []string{}
	thumbnailImagesCollection := []string{}
	files := p.Files

	for i := range files {
		modelFile := strategies.UploadedImageModel{
			Bytes:    files[i].Bytes,
			Filename: files[i].Name,
			Type:     files[i].Type,
		}
		filenameToSave := uploadImage(modelFile, s.strategy)
		thumbnailFilenameToSave := processThumbnail(modelFile, s.strategy, s.config, filenameToSave, s.logger)

		//s.config.ThumbnailMaxHeight

		imagesCollection = append(imagesCollection, filenameToSave)
		thumbnailImagesCollection = append(thumbnailImagesCollection, thumbnailFilenameToSave)
	}

	total := len(imagesCollection)

	// And return it
	res = &images.ImagesListMedia{
		Images:          imagesCollection,
		ThumbnailImages: thumbnailImagesCollection,
		Total:           &total,
	}

	return
}

func processThumbnail(modelFile strategies.UploadedImageModel, strategy strategies.StorageStrategy, config *config.ServiceConfig, regularImageFilename string, logger *log.Logger) string { //data []byte, imageType string, config *config.ServiceConfig) (outputFilename string) {
	data := modelFile.Bytes
	if uint(len(data)) > config.ThumbnailMaxSize && (*modelFile.Type == "image/jpeg" || *modelFile.Type == "image/png") {
		originalImage, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			// have a look at this issue for a solution to some of the errors caught here (invalid JPEG format: short Huffman data)
			// https://github.com/golang/go/issues/10447
			logger.Printf("Error decoding image to create thumbsnail image: %s", err.Error())
			return regularImageFilename
		}
		newImage := resize.Thumbnail(config.ThumbnailMaxWidth, config.ThumbnailMaxHeight, originalImage, resize.Lanczos3)
		buf := new(bytes.Buffer)
		if *modelFile.Type == "image/jpeg" {
			options := &jpeg.Options{Quality: int(config.ThumbnailQuality)}
			err = jpeg.Encode(buf, newImage, options)
		} else {

			enc := &png.Encoder{
				CompressionLevel: mapCode[config.ThumbnailCompression],
			}
			err = enc.Encode(buf, newImage)
		}
		if err != nil {
			logger.Printf("Error encoding image to create thumbsnail image: %s", err)
			return regularImageFilename
		}
		thumbFilename := "thumb-" + *modelFile.Filename
		modelFile.Filename = &thumbFilename
		modelFile.Bytes = buf.Bytes()
		thumbnailFilenameToSave := uploadImage(modelFile, strategy)
		return thumbnailFilenameToSave
	}
	return regularImageFilename
}

func uploadImage(modelFile strategies.UploadedImageModel, strategy strategies.StorageStrategy) string {
	filenameToSave, err := strategy.UploadImage(&modelFile)
	if err != nil {
		panic(err)
	}

	firstChar := filenameToSave[0]
	if string(firstChar) == "." {
		filenameToSave = filenameToSave[1:]
	}
	return filenameToSave
}
