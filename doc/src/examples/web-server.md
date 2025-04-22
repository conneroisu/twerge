# Web Server Integration

This page demonstrates how to integrate Twerge with a Go web server for delivering optimized Tailwind CSS in a production environment.

## HTTP Server Example

The example below shows how to set up a simple HTTP server that serves a web application with Twerge-optimized Tailwind CSS classes.

### Server Setup

```go title="main.go"
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/conneroisu/twerge/examples/simple/views"
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("_static"))
	http.Handle("/dist/", http.StripPrefix("/dist/", fs))
	
	// Index route
	http.HandleFunc("/", handleIndex)
	
	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// If not at the root path, return 404
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	// Render the template
	err := views.View().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper to get environment variable with default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
```

### Middleware for Cache Control

For production, you might want to add cache control headers for better performance:

```go title="middleware.go"
package main

import (
	"net/http"
	"strings"
	"time"
)

// CacheControlMiddleware adds cache control headers for static assets
func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set cache headers for static assets
		if strings.HasPrefix(r.URL.Path, "/dist/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
			w.Header().Set("Expires", time.Now().Add(time.Hour*24*365).Format(time.RFC1123))
		}
		next.ServeHTTP(w, r)
	})
}

// Usage example
func setupServer() {
	staticHandler := http.FileServer(http.Dir("_static"))
	http.Handle("/dist/", CacheControlMiddleware(http.StripPrefix("/dist/", staticHandler)))
}
```

## Template Composition

When building larger applications, you'll want to compose templates. Here's how to use Twerge with template composition:

```go title="layout.templ"
package views

import "github.com/conneroisu/twerge"

templ Layout(title string) {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<title>{ title }</title>
			<link rel="stylesheet" href="/dist/styles.css"/>
		</head>
		<body class={ twerge.It("bg-gray-50 text-gray-900 flex flex-col min-h-screen") }>
			<header class={ twerge.It("bg-indigo-600 text-white shadow-md") }>
				<!-- Header content -->
			</header>
			
			<!-- Slot for page content -->
			{ children... }
			
			<footer class={ twerge.It("bg-gray-800 text-white py-6") }>
				<!-- Footer content -->
			</footer>
		</body>
	</html>
}

// Usage in page templates
templ AboutPage() {
	@Layout("About Us") {
		<main class={ twerge.It("container mx-auto px-4 py-6 flex-grow") }>
			<h1 class={ twerge.It("text-3xl font-bold mb-6") }>About Us</h1>
			<p class={ twerge.It("text-gray-700") }>Content goes here...</p>
		</main>
	}
}
```

## API Routes with JSON

If you're building an API that serves both HTML and JSON, you can structure your handlers like this:

```go title="api_handler.go"
package main

import (
	"encoding/json"
	"net/http"
	
	"github.com/conneroisu/twerge/examples/simple/views"
)

// Struct for JSON response
type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Handler that can respond with HTML or JSON
func handleData(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"name": "Twerge Example",
		"description": "A demonstration of Twerge with web servers",
	}
	
	// Check Accept header for JSON request
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		response := ApiResponse{
			Status: "success",
			Data: data,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Otherwise render HTML
	err := views.DataView(data).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
```

## Performance Considerations

When using Twerge with a web server, consider these performance optimizations:

1. **Precompiled Templates**: Generate all templ components at build time
2. **Static File Serving**: Use proper cache headers for static assets
3. **Compression**: Enable gzip/brotli compression for CSS and HTML
4. **CDN Integration**: Consider serving static assets from a CDN

```go title="compression.go"
package main

import (
	"net/http"
	"strings"
	
	"github.com/NYTimes/gziphandler"
)

func setupCompression() {
	// Apply gzip compression to static assets
	staticHandler := http.FileServer(http.Dir("_static"))
	compressedHandler := gziphandler.GzipHandler(http.StripPrefix("/dist/", staticHandler))
	http.Handle("/dist/", compressedHandler)
}
```

## Example Server with All Features

Here's a complete example integrating all the features mentioned above:

```go title="main.go"
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/NYTimes/gziphandler"
	"github.com/conneroisu/twerge/examples/simple/views"
)

func main() {
	// Create router
	mux := http.NewServeMux()
	
	// Static files with compression and caching
	staticHandler := http.FileServer(http.Dir("_static"))
	compressedHandler := gziphandler.GzipHandler(http.StripPrefix("/dist/", staticHandler))
	cachedHandler := CacheControlMiddleware(compressedHandler)
	mux.Handle("/dist/", cachedHandler)
	
	// Routes
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/data", handleData)
	
	// Configure server
	server := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on http://localhost%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server gracefully stopped")
}
```

This comprehensive example demonstrates how to integrate Twerge with a production-ready web server, complete with compression, caching, and graceful shutdown.