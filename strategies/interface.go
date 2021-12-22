package strategies

// StorageStrategy ...
type StorageStrategy interface {
	UploadImage(uploadedImage *UploadedImageModel) (string, error)
}
