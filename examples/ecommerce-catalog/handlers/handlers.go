package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/conneroisu/twerge/examples/ecommerce-catalog/data"
	"github.com/conneroisu/twerge/examples/ecommerce-catalog/types"
	"github.com/conneroisu/twerge/examples/ecommerce-catalog/views"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	store *data.Store
}

func New(store *data.Store) *Handlers {
	return &Handlers{store: store}
}

// CatalogPage renders the main catalog page
func (h *Handlers) CatalogPage(w http.ResponseWriter, r *http.Request) {
	filters := h.parseFilters(r)
	products := h.store.GetProducts(filters)
	categoryCounts := h.store.GetCategoryCounts()
	cartItems := h.store.GetCartItemMap("default-session")
	favorites := h.store.GetFavorites("default-session")
	cartCount := h.store.GetCartCount("default-session")

	component := views.CatalogPage(
		products,
		filters,
		len(products),
		categoryCounts,
		cartItems,
		favorites,
		cartCount,
	)

	component.Render(r.Context(), w)
}

// ProductPage renders individual product detail page
func (h *Handlers) ProductPage(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	product, err := h.store.GetProductByID(productID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	selectedVariants := make(map[string]string)
	cartItems := h.store.GetCartItemMap("default-session")
	cartCount := h.store.GetCartCount("default-session")

	component := views.ProductDetail(product, selectedVariants, cartItems, cartCount)
	component.Render(r.Context(), w)
}

// GetProducts returns filtered products (HTMX endpoint)
func (h *Handlers) GetProducts(w http.ResponseWriter, r *http.Request) {
	filters := h.parseFilters(r)
	products := h.store.GetProducts(filters)
	cartItems := h.store.GetCartItemMap("default-session")
	favorites := h.store.GetFavorites("default-session")

	component := views.CatalogContent(
		products,
		filters,
		len(products),
		cartItems,
		favorites,
	)

	component.Render(r.Context(), w)
}

// FilterProducts handles filter form submissions
func (h *Handlers) FilterProducts(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	filters := h.parseFilters(r)
	products := h.store.GetProducts(filters)
	cartItems := h.store.GetCartItemMap("default-session")
	favorites := h.store.GetFavorites("default-session")

	component := views.CatalogContent(
		products,
		filters,
		len(products),
		cartItems,
		favorites,
	)

	component.Render(r.Context(), w)
}

// AddToCart handles adding items to cart
func (h *Handlers) AddToCart(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")
	
	product, err := h.store.GetProductByID(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if !product.InStock {
		component := views.Toast("Product is out of stock", "error")
		component.Render(r.Context(), w)
		return
	}

	err = h.store.AddToCart("default-session", productID, 1, []string{})
	if err != nil {
		component := views.Toast("Failed to add item to cart", "error")
		component.Render(r.Context(), w)
		return
	}

	component := views.Toast(fmt.Sprintf("Added %s to cart", product.Name), "success")
	component.Render(r.Context(), w)
}

// RemoveFromCart handles removing items from cart
func (h *Handlers) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")
	
	err := h.store.RemoveFromCart("default-session", itemID)
	if err != nil {
		component := views.Toast("Failed to remove item from cart", "error")
		component.Render(r.Context(), w)
		return
	}

	// Return empty content to remove the item from the DOM
	w.WriteHeader(http.StatusOK)
}

// UpdateCartItem handles quantity updates
func (h *Handlers) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	quantityStr := r.FormValue("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	if quantity == 0 {
		// Remove item if quantity is 0
		h.RemoveFromCart(w, r)
		return
	}

	item, err := h.store.UpdateCartItemQuantity("default-session", itemID, quantity)
	if err != nil {
		component := views.Toast("Failed to update cart item", "error")
		component.Render(r.Context(), w)
		return
	}

	component := views.CartItem(item, false)
	component.Render(r.Context(), w)
}

// GetCart returns the shopping cart drawer
func (h *Handlers) GetCart(w http.ResponseWriter, r *http.Request) {
	cart := h.store.GetCart("default-session")
	component := views.CartDrawer(cart, true)
	component.Render(r.Context(), w)
}

// GetCartCount returns just the cart count for the badge
func (h *Handlers) GetCartCount(w http.ResponseWriter, r *http.Request) {
	count := h.store.GetCartCount("default-session")
	component := views.CartBadge(count)
	component.Render(r.Context(), w)
}

// SelectVariant handles product variant selection
func (h *Handlers) SelectVariant(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	variantType := r.FormValue("type")
	variantValue := r.FormValue("value")

	product, err := h.store.GetProductByID(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Update selected variants (in a real app, this would be stored in session)
	selectedVariants := map[string]string{
		variantType: variantValue,
	}

	// Find the variant to check availability
	var selectedVariant *types.ProductVariant
	for _, variant := range product.Variants {
		if variant.Type == variantType && variant.Value == variantValue {
			selectedVariant = &variant
			break
		}
	}

	if selectedVariant != nil {
		component := views.VariantAvailabilityAlert(variantType, variantValue, selectedVariant.Available)
		component.Render(r.Context(), w)
	} else {
		// Return updated variant selectors
		component := views.ProductVariantSelectors(product, selectedVariants)
		component.Render(r.Context(), w)
	}
}

// parseFilters extracts filter parameters from the request
func (h *Handlers) parseFilters(r *http.Request) types.FilterState {
	filters := types.FilterState{
		Categories:  []string{},
		PriceRange:  [2]float64{0, 1000},
		Ratings:     []int{},
		SearchQuery: r.URL.Query().Get("search"),
		SortBy:      r.URL.Query().Get("sort"),
	}

	// Parse form if it's a POST request
	if r.Method == "POST" {
		r.ParseForm()
	}

	// Categories
	categories := r.Form["categories"]
	if len(categories) == 0 {
		categories = r.URL.Query()["category"]
	}
	filters.Categories = categories

	// Price range
	if minPrice := r.FormValue("price_min"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters.PriceRange[0] = price
		}
	}
	if maxPrice := r.FormValue("price_max"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters.PriceRange[1] = price
		}
	}

	// Availability filters
	if r.FormValue("in_stock") == "true" {
		inStock := true
		filters.InStock = &inStock
	}
	if r.FormValue("on_sale") == "true" {
		onSale := true
		filters.OnSale = &onSale
	}

	// Ratings
	ratingStrs := r.Form["ratings"]
	for _, ratingStr := range ratingStrs {
		if rating, err := strconv.Atoi(ratingStr); err == nil {
			filters.Ratings = append(filters.Ratings, rating)
		}
	}

	// Search query
	if search := r.FormValue("search"); search != "" {
		filters.SearchQuery = search
	}

	// Sort
	if sortBy := r.FormValue("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}
	if filters.SortBy == "" {
		filters.SortBy = "relevance"
	}

	return filters
}

// JSON API endpoints for development/testing

// GetProductsAPI returns products as JSON
func (h *Handlers) GetProductsAPI(w http.ResponseWriter, r *http.Request) {
	filters := h.parseFilters(r)
	products := h.store.GetProducts(filters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"total":    len(products),
		"filters":  filters,
	})
}

// GetCartAPI returns cart as JSON
func (h *Handlers) GetCartAPI(w http.ResponseWriter, r *http.Request) {
	cart := h.store.GetCart("default-session")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// Helper functions

func parseCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}