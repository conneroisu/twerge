package types

import "time"

// Product represents a product in the catalog
type Product struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	SalePrice   *float64         `json:"sale_price,omitempty"`
	InStock     bool             `json:"in_stock"`
	Featured    bool             `json:"featured"`
	Images      []string         `json:"images"`
	Variants    []ProductVariant `json:"variants"`
	Category    string           `json:"category"`
	Rating      float64          `json:"rating"`
	Reviews     int              `json:"reviews"`
	CreatedAt   time.Time        `json:"created_at"`
}

// ProductVariant represents a variant of a product (size, color, style)
type ProductVariant struct {
	Type      string   `json:"type"`      // "color", "size", "style"
	Value     string   `json:"value"`     // "red", "medium", "classic"
	Available bool     `json:"available"`
	Price     *float64 `json:"price,omitempty"` // nil if same as base price
	SKU       string   `json:"sku,omitempty"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID               string    `json:"id"`
	Product          Product   `json:"product"`
	Quantity         int       `json:"quantity"`
	SelectedVariants []string  `json:"selected_variants"`
	Price            float64   `json:"price"` // Final price including variant prices
	AddedAt          time.Time `json:"added_at"`
}

// Cart represents a shopping cart
type Cart struct {
	ID        string     `json:"id"`
	Items     []CartItem `json:"items"`
	Subtotal  float64    `json:"subtotal"`
	Tax       float64    `json:"tax"`
	Shipping  float64    `json:"shipping"`
	Total     float64    `json:"total"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// FilterState represents the current filter state
type FilterState struct {
	Categories   []string    `json:"categories"`
	PriceRange   [2]float64  `json:"price_range"`
	InStock      *bool       `json:"in_stock,omitempty"`
	OnSale       *bool       `json:"on_sale,omitempty"`
	Ratings      []int       `json:"ratings"`
	SearchQuery  string      `json:"search_query"`
	SortBy       string      `json:"sort_by"` // "price_asc", "price_desc", "rating", "newest"
}

// Categories available in the catalog
var Categories = []string{
	"Electronics",
	"Clothing",
	"Home & Garden",
	"Sports & Outdoors",
	"Books",
	"Toys & Games",
	"Beauty & Personal Care",
	"Automotive",
}

// Color mapping for color variants
var ColorMap = map[string]string{
	"red":     "bg-red-500",
	"blue":    "bg-blue-500",
	"green":   "bg-green-500",
	"yellow":  "bg-yellow-500",
	"purple":  "bg-purple-500",
	"pink":    "bg-pink-500",
	"black":   "bg-gray-900",
	"white":   "bg-white border-2 border-gray-300",
	"gray":    "bg-gray-500",
	"brown":   "bg-yellow-800",
	"orange":  "bg-orange-500",
	"teal":    "bg-teal-500",
	"indigo":  "bg-indigo-500",
}

// Size options
var Sizes = []string{"XS", "S", "M", "L", "XL", "XXL"}

// Style options
var Styles = []string{"Classic", "Modern", "Vintage", "Sport", "Casual", "Formal"}