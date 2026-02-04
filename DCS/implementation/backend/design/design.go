package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl" // Kein Punkt, sondern Alias 'cors'
)

var _ = API("dcs", func() {
	Title("DCS API Server")
	Version("0.0.1")

	cors.Origin("*", func() {
		cors.Headers("X-Shared-Secret", "X-Api-Version")
		cors.MaxAge(100)
		cors.Credentials()
	})
	Server("dcs", func() {
		Host("local", func() {
			URI("http://0.0.0.0:8991")
		})
	})
})
