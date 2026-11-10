package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	publicDir := "public"
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		publicDir = "."
	}

	fileServer := http.FileServer(http.Dir(publicDir))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Proper MIME type for WebAssembly
		if strings.HasSuffix(path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}

		// Prevent aggressive caching during game jam / testing
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		// Fallback to index.html if route is root
		if path == "/" || path == "" {
			http.ServeFile(w, r, filepath.Join(publicDir, "index.html"))
			return
		}

		fileServer.ServeHTTP(w, r)
	})

	log.Printf("==================================================")
	log.Printf("  🏛️  ATHENS RAMPAGE: ΤΙΤΑΝΕΣ ΣΤΗΝ ΑΚΡΟΠΟΛΗ 🏛️  ")
	log.Printf("  Serving WebAssembly on http://localhost:%s", port)
	log.Printf("==================================================")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
