package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/conneroisu/twerge/examples/admin-dashboard/views"
)

func main() {
	// Initialize the classes cache
	// This will be populated by the generated code
	if _, err := os.Stat("classes/classes.go"); err == nil {
		// Import will be handled automatically when generated
		log.Println("Classes cache initialized")
	}

	// Set up routes
	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/static/", handleStatic)

	// Start server
	port := ":8081"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	// Get sidebar state from query params
	collapsed := r.URL.Query().Get("sidebar") == "collapsed"
	
	// Render the dashboard
	err := views.Dashboard(collapsed).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	// Serve static files from the static directory
	path := r.URL.Path[len("/static/"):]
	
	// Determine the file path
	var filePath string
	switch path {
	case "styles.css":
		filePath = filepath.Join("static", "dist", "styles.css")
	default:
		filePath = filepath.Join("static", path)
	}
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	
	// Serve the file
	http.ServeFile(w, r, filePath)
}