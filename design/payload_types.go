package design

import (
	. "goa.design/goa/v3/dsl"
)

// ImageUpload single image upload element
var ImageUpload = Type("ImageUpload", func() {
	Description("A single Image Upload type")
	Attribute("type", String)
	Attribute("bytes", Bytes)
	Attribute("name", String)
})

// ImageUploadPayload is a list of files
var ImageUploadPayload = Type("ImageUploadPayload", func() {
	Description("Image Upload Payload")

	Attribute("Files", ArrayOf(ImageUpload), "Collection of uploaded files")
})
