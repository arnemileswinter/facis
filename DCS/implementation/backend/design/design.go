package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl" // Kein Punkt, sondern Alias 'cors'
)

var _ = API("dcs", func() {
	Title("DCS API Server")
	Version("0.0.1")

	cors.Origin("*", func() {
		cors.Headers("Content-Type", "Authorization", "X-Shared-Secret", "X-Api-Version")
		cors.Methods("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")
		cors.MaxAge(100)
	})
	Server("dcs", func() {
		Host("local", func() {
			URI("http://0.0.0.0:8991")
		})
	})
})
