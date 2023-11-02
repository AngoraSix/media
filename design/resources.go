package design // The convention consists of naming the design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("images", func() { // Resources group related API endpoints
	HTTP(func() { // together. They map to REST resources for REST services.
		Path("/images")
	})

	Method("upload", func() { // Actions define a single API endpoint together with its path, parameters (both pathparameters and querystring values) and payload Responses define the shape and status code
		Description("Uploads a new image")
		HTTP(func() {
			POST("/")
			MultipartRequest()
			Response("InternalServerError", StatusInternalServerError)
		})

		Payload(ImageUploadPayload)

		Result(ImageList)
		Error("InternalServerError")
	})
})

var _ = Service("static", func() { // Resources group related API endpoints
	cors.Origin("*")

	Files("/static/uploads/{filename}", "./static/uploads")
})
