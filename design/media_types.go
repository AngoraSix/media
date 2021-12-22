package design

import (
	. "goa.design/goa/v3/dsl"
)

// ImageList is a list of images
var ImageList = ResultType("application/vnd.images.list.media+json", func() {
	Description("A list of images")

	Attributes(func() {
		Attribute("total", Int, "Total companies found.", func() {})

		Attribute("images", ArrayOf(String), "A list of images url.")
		Attribute("thumbnailImages", ArrayOf(String), "A list of thumbail images url (light version of the images collection).")
	})

	View("default", func() {
		Attribute("total")
		Attribute("images")
		Attribute("thumbnailImages")
	})
})
