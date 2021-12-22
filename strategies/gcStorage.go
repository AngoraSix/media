package strategies

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"angorasix.com/media/config"
	"cloud.google.com/go/storage"
)

// GCPBucketCreationErrors ...
var GCPBucketCreationErrors = []string{
	"already own this bucket",
}

// GoogleCloudStrategy ...
type GoogleCloudStrategy struct {
	client            *storage.Client
	bucket            *storage.BucketHandle
	gpcStorageApiHost string
}

// CreateGoogleCloudStrategy ...
func CreateGoogleCloudStrategy(config *config.ServiceConfig) (StorageStrategy, error) {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := config.ProjectID

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(config.BucketName)

	// First we check if bucket exists...
	_, err = bucket.Attrs(ctx)
	if err != nil {
		// Creates the new bucket.
		err = bucket.Create(ctx, projectID, nil)
		if err != nil {
			creatingError := true

			for _, allowedErrorString := range GCPBucketCreationErrors {
				if strings.Contains(strings.ToLower(err.Error()), strings.ToLower(allowedErrorString)) {
					creatingError = false
					break
				}
			}

			if !creatingError {
				fmt.Printf("Creation of Bucket returns: %s \n\n", err.Error())
			} else {
				return nil, err
			}
		} else {
			fmt.Printf("Bucket %s created.\n", config.BucketName)
		}
	}

	strategy := &GoogleCloudStrategy{
		client,
		bucket,
		config.StorageAPIHost,
	}

	return strategy, nil
}

// UploadImage uploads image
func (s *GoogleCloudStrategy) UploadImage(img *UploadedImageModel) (string, error) {
	// Creates filename.
	filenameToSave := fmt.Sprintf("%s_%s", createNowString(), *img.Filename)

	// Creates Google Storage Object and Writer
	ctx := context.Background()
	object := s.bucket.Object(filenameToSave)
	wc := object.NewWriter(ctx)
	wc.ContentType = *img.Type
	wc.CacheControl = "public, max-age=350000"
	if _, err := io.Copy(wc, bytes.NewReader(img.Bytes)); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	gcObjectAttrs := wc.Attrs()
	url := fmt.Sprintf("%s/%s/%s", s.gpcStorageApiHost, gcObjectAttrs.Bucket, gcObjectAttrs.Name)

	return url, nil
}
