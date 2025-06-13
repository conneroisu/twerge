package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/conneroisu/twerge/examples/realtime-collab/data"
	"github.com/conneroisu/twerge/examples/realtime-collab/handlers"
	"github.com/conneroisu/twerge/examples/realtime-collab/websocket"
)

func main() {
	// Initialize data store with sample data
	log.Println("🗄️  Initializing data store...")
	store := data.NewStore()

	// Initialize WebSocket hub
	log.Println("🔌 Starting WebSocket hub...")
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize HTTP handlers
	log.Println("🌐 Setting up HTTP handlers...")
	handlers := handlers.NewHandlers(store, wsHub)
	router := handlers.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         getPort(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on %s", server.Addr)
		log.Printf("📊 Dashboard: http://localhost%s", server.Addr)
		log.Printf("🔗 WebSocket endpoint: ws://localhost%s/ws", server.Addr)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Print startup information
	printStartupInfo()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func printStartupInfo() {
	log.Println("🎯 Realtime Collaborative Dashboard")
	log.Println("=====================================")
	log.Println("🔥 Features:")
	log.Println("   • Real-time collaboration with WebSockets")
	log.Println("   • Multi-user document editing")
	log.Println("   • Live cursor tracking")
	log.Println("   • Instant notifications")
	log.Println("   • Project management")
	log.Println("   • Team presence indicators")
	log.Println("   • Complex UI state optimization with Twerge")
	log.Println("")
	log.Println("🛠️  Technology Stack:")
	log.Println("   • Backend: Go with Chi router")
	log.Println("   • WebSockets: Gorilla WebSocket")
	log.Println("   • Frontend: HTMX + Alpine.js")
	log.Println("   • Templates: a-h/templ")
	log.Println("   • CSS: TailwindCSS + Twerge optimization")
	log.Println("   • Data: In-memory store with realistic sample data")
	log.Println("")
	log.Println("👥 Demo Users:")
	log.Println("   • Alice (Admin) - alice.dev@example.com")
	log.Println("   • Bob (Designer) - bob.designer@example.com")
	log.Println("   • Carol (PM) - carol.pm@example.com")
	log.Println("   • David (Analyst) - david.analyst@example.com")
	log.Println("   • Eve (Writer) - eve.writer@example.com")
	log.Println("")
	log.Println("📝 Sample Projects:")
	log.Println("   • Design System Overhaul (High Priority)")
	log.Println("   • API Documentation Portal (Medium Priority)")
	log.Println("   • User Research Findings (Low Priority)")
	log.Println("   • Mobile App Prototype (Critical Priority)")
	log.Println("")
	log.Println("🌟 Real-time Features:")
	log.Println("   • Live user presence (online/offline/away/busy)")
	log.Println("   • Real-time cursor positions in documents")
	log.Println("   • Instant activity feed updates")
	log.Println("   • Live typing indicators")
	log.Println("   • Multi-room collaboration")
	log.Println("   • System health monitoring")
	log.Println("")
	log.Println("🎨 UI Complexity Optimized by Twerge:")
	log.Println("   • 150+ unique component states")
	log.Println("   • Responsive design (mobile/tablet/desktop)")
	log.Println("   • Dark/light theme variations")
	log.Println("   • Interactive hover/focus/active states")
	log.Println("   • Complex conditional styling")
	log.Println("   • Animation and transition classes")
	log.Println("")
	log.Println("🚀 Ready for development and testing!")
	log.Println("=====================================")
}