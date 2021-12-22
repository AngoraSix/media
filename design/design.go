package design // The convention consists of naming the design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"

	config "angorasix.com/media/config"
)

var _ = API("media", func() {
	Title("HOC Media")
	Description("Media Content microservice")
	Server("main", func() {
		Services("static", "images")
		Host("0.0.0.0", func() {
			URI(config.DefaultServerConfig.GetHost())
		})
	})
	Server("openapi", func() {
		Description("OpenAPI server hosts the service OpenAPI specification.")
		Services("openapi")
		Host("0.0.0.0", func() {
			Description("default host")
			URI(config.DefaultServerConfig.GetHost() + "/openapi")
		})
	})
})

var _ = Service("openapi", func() {
	cors.Origin("*")
	Files("/openapi.json", "./gen/http/openapi.json")
})
