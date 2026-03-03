package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	goahttp "goa.design/goa/v3/http"
)

// mountFrontend registers the frontend static file server routes on the given Goa mux.
// It serves the built Vue.js frontend at /ui and redirects root to /ui.
func mountFrontend(mux goahttp.Muxer) {
	const staticDir = "/app/web/dist"

	// Check if the static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// Frontend not embedded, skip mounting
		return
	}

	// Redirect root to /ui
	mux.Handle("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusMovedPermanently)
	})

	// Serve frontend at /ui with SPA routing
	mux.Handle("GET", "/ui", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/", http.StatusMovedPermanently)
	})

	mux.Handle("GET", "/ui/*", func(w http.ResponseWriter, r *http.Request) {
		// Strip /ui prefix to get the actual file path
		path := strings.TrimPrefix(r.URL.Path, "/ui")
		if path == "" {
			path = "/"
		}

		// Clean the path to prevent directory traversal
		path = filepath.Clean(path)

		// Build the full path
		fullPath := filepath.Join(staticDir, path)

		// Resolve to absolute path and verify it's within staticDir
		absStaticDir, err := filepath.Abs(staticDir)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		absFullPath, err := filepath.Abs(fullPath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Verify the resolved path is still within the static directory
		relPath, err := filepath.Rel(absStaticDir, absFullPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			http.NotFound(w, r)
			return
		}

		// Try to serve the file
		if info, err := os.Stat(absFullPath); err == nil && !info.IsDir() {
			// File exists, serve it
			http.ServeFile(w, r, absFullPath)
			return
		}

		// File doesn't exist or path is a directory - serve index.html for SPA routing
		indexPath := filepath.Join(absStaticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(w, r, indexPath)
			return
		}

		// index.html not found
		http.NotFound(w, r)
	})
}
