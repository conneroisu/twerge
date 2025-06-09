package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/conneroisu/twerge/examples/ecommerce-catalog/data"
	"github.com/conneroisu/twerge/examples/ecommerce-catalog/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed _static/dist/*
var staticFiles embed.FS

func main() {
	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// Initialize handlers
	h := handlers.New(data.NewStore())

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	// Routes
	r.Get("/", h.CatalogPage)
	r.Get("/products/{id}", h.ProductPage)
	
	// HTMX endpoints
	r.Post("/cart/add/{productId}", h.AddToCart)
	r.Delete("/cart/remove/{itemId}", h.RemoveFromCart)
	r.Patch("/cart/update/{itemId}", h.UpdateCartItem)
	r.Get("/cart", h.GetCart)
	r.Get("/cart/count", h.GetCartCount)
	
	// Filter endpoints
	r.Post("/catalog/filter", h.FilterProducts)
	r.Get("/catalog/products", h.GetProducts)
	
	// Variant selection
	r.Post("/products/{id}/variants", h.SelectVariant)

	// Server configuration
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		fmt.Println("🛍️  E-commerce Catalog Server starting on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")
	srv.Close()
}