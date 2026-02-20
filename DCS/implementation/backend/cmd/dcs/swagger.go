package main

import (
	"net/http"
	"os"

	goahttp "goa.design/goa/v3/http"
)

// mountSwaggerUI registers Swagger UI routes on the given Goa mux:
//
//	GET /swagger      → Swagger UI (loaded from CDN)
//	GET /openapi3.json → Generated OpenAPI 3 specification
func mountSwaggerUI(mux goahttp.Muxer) {
	// Serve the Swagger UI HTML page.
	mux.Handle("GET", "/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(swaggerHTML))
	})
	mux.Handle("GET", "/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger", http.StatusMovedPermanently)
	})

	// Serve the generated OpenAPI 3 specification JSON.
	mux.Handle("GET", "/openapi3.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("gen/http/openapi3.json")
		if err != nil {
			http.Error(w, "OpenAPI spec not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data)
	})
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>DCS API – Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/openapi3.json",
      dom_id: "#swagger-ui",
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.SwaggerUIStandalonePreset,
      ],
      layout: "BaseLayout",
      deepLinking: true,
    });
  </script>
</body>
</html>`
