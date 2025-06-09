//go:build ignore

package main

import (
	"log"
	"time"

	"github.com/conneroisu/twerge"
	"github.com/conneroisu/twerge/examples/ecommerce-catalog/types"
	"github.com/conneroisu/twerge/examples/ecommerce-catalog/views"
)

func main() {
	// Create sample data for code generation
	sampleProduct := types.Product{
		ID:          "sample-001",
		Name:        "Sample Product for Code Generation",
		Description: "This is a sample product used for generating optimized CSS classes with Twerge.",
		Price:       99.99,
		SalePrice:   ptr(79.99),
		InStock:     true,
		Featured:    true,
		Images:      []string{"/static/images/sample-1.jpg", "/static/images/sample-2.jpg"},
		Category:    "Electronics",
		Rating:      4.5,
		Reviews:     123,
		CreatedAt:   time.Now(),
		Variants: []types.ProductVariant{
			{Type: "color", Value: "black", Available: true},
			{Type: "color", Value: "white", Available: true},
			{Type: "color", Value: "blue", Available: false},
			{Type: "size", Value: "S", Available: true},
			{Type: "size", Value: "M", Available: true},
			{Type: "size", Value: "L", Available: true},
			{Type: "style", Value: "Standard", Available: true},
			{Type: "style", Value: "Premium", Available: true, Price: ptr(20.0)},
		},
	}

	outOfStockProduct := types.Product{
		ID:          "sample-002",
		Name:        "Out of Stock Sample Product",
		Description: "This product is out of stock for testing purposes.",
		Price:       49.99,
		InStock:     false,
		Featured:    false,
		Images:      []string{"/static/images/sample-3.jpg"},
		Category:    "Clothing",
		Rating:      4.0,
		Reviews:     56,
		CreatedAt:   time.Now().AddDate(0, -1, 0),
		Variants: []types.ProductVariant{
			{Type: "color", Value: "red", Available: false},
			{Type: "size", Value: "M", Available: false},
		},
	}

	sampleCartItem := types.CartItem{
		ID:      "cart-001",
		Product: sampleProduct,
		Quantity: 2,
		SelectedVariants: []string{"Black", "Medium"},
		Price:   79.99,
		AddedAt: time.Now(),
	}

	sampleCart := types.Cart{
		ID:       "sample-cart",
		Items:    []types.CartItem{sampleCartItem},
		Subtotal: 159.98,
		Tax:      12.80,
		Shipping: 0.00,
		Total:    172.78,
		UpdatedAt: time.Now(),
	}

	sampleFilters := types.FilterState{
		Categories:  []string{"Electronics", "Clothing"},
		PriceRange:  [2]float64{0, 200},
		InStock:     ptr(true),
		OnSale:      ptr(true),
		Ratings:     []int{4, 5},
		SearchQuery: "sample",
		SortBy:      "price_asc",
	}

	sampleProducts := []types.Product{sampleProduct, outOfStockProduct}
	cartItems := map[string]bool{
		"sample-001": true,
		"sample-002": false,
	}
	favorites := map[string]bool{
		"sample-001": true,
		"sample-002": false,
	}
	categoryCounts := map[string]int{
		"Electronics": 15,
		"Clothing":    23,
		"Home & Garden": 18,
		"Sports & Outdoors": 12,
		"Books": 8,
		"Toys & Games": 6,
		"Beauty & Personal Care": 14,
		"Automotive": 9,
	}
	selectedVariants := map[string]string{
		"color": "black",
		"size":  "M",
		"style": "Premium",
	}

	log.Println("🎨 Generating optimized CSS classes with Twerge...")

	// Generate all component variations for comprehensive class coverage
	if err := twerge.CodeGen(
		twerge.Default(),
		"classes/classes.go",
		"input.css",
		"classes/classes.html",

		// Layout components - all states
		views.Layout("Sample Page Title", 0),
		views.Layout("Sample Page Title", 3),
		views.Header(0),
		views.Header(5),
		views.Footer(),
		views.Toast("Success message", "success"),
		views.Toast("Error message", "error"),
		views.Toast("Info message", "info"),

		// Product components - all variations
		views.ProductCard(sampleProduct, false, false),   // Normal state
		views.ProductCard(sampleProduct, true, false),    // In cart
		views.ProductCard(sampleProduct, false, true),    // Favorited
		views.ProductCard(sampleProduct, true, true),     // In cart + favorited
		views.ProductCard(outOfStockProduct, false, false), // Out of stock
		views.ProductGrid(sampleProducts, cartItems, favorites),

		// Product detail page
		views.ProductDetail(sampleProduct, selectedVariants, cartItems, 3),

		// Variant selectors - all types
		views.VariantSelector(sampleProduct.Variants, "black", "color", "sample-001"),
		views.VariantSelector(sampleProduct.Variants, "M", "size", "sample-001"),
		views.VariantSelector(sampleProduct.Variants, "Premium", "style", "sample-001"),
		views.VariantOption(sampleProduct.Variants[0], true, "color", "sample-001"),   // Selected color
		views.VariantOption(sampleProduct.Variants[0], false, "color", "sample-001"),  // Unselected color
		views.VariantOption(sampleProduct.Variants[2], false, "color", "sample-001"),  // Unavailable color
		views.ProductVariantSelectors(sampleProduct, selectedVariants),
		views.SelectedVariantsSummary(selectedVariants, sampleProduct.Variants),
		views.VariantAvailabilityAlert("color", "blue", false),
		views.VariantAvailabilityAlert("size", "M", true),
		views.QuickVariantSelector(sampleProduct.Variants, "color", "sample-001"),
		views.VariantsPreview(sampleProduct.Variants),

		// Cart components - all states
		views.CartDrawer(sampleCart, true),
		views.CartDrawer(types.Cart{Items: []types.CartItem{}}, true), // Empty cart
		views.CartItem(sampleCartItem, false),  // Normal state
		views.CartItem(sampleCartItem, true),   // Updating state
		views.EmptyCart(),
		views.CartFooter(sampleCart),
		views.CartBadge(0),
		views.CartBadge(3),
		views.CartBadge(99),
		views.AddToCartButton(sampleProduct, false, false),    // Available, not in cart
		views.AddToCartButton(sampleProduct, true, false),     // Available, in cart
		views.AddToCartButton(outOfStockProduct, false, false), // Out of stock
		views.AddToCartButton(sampleProduct, false, true),     // Disabled
		views.QuickAddButton("sample-001", "small"),
		views.QuickAddButton("sample-001", "large"),

		// Filter components - all states
		views.FilterSidebar(sampleFilters, 42, categoryCounts),
		views.CategoryFilter([]string{"Electronics"}, categoryCounts),
		views.CategoryFilter([]string{}, categoryCounts), // No selection
		views.FilterCheckbox("Electronics", true, 15, "categories"),
		views.FilterCheckbox("Electronics", false, 15, "categories"),
		views.FilterCheckbox("Electronics", false, 0, "categories"), // No items
		views.PriceRangeFilter([2]float64{25, 100}),
		views.PriceRangeFilter([2]float64{0, 1000}), // Full range
		views.PriceQuickFilter("Under $25", [2]float64{0, 25}, [2]float64{0, 25}),   // Selected
		views.PriceQuickFilter("Under $25", [2]float64{0, 25}, [2]float64{50, 100}), // Not selected
		views.AvailabilityFilter(ptr(true), ptr(true)),
		views.AvailabilityFilter(ptr(false), ptr(false)),
		views.AvailabilityFilter(nil, nil), // No filters
		views.RatingFilter([]int{4, 5}),
		views.RatingFilter([]int{}), // No selection
		views.RatingOption(5, true),
		views.RatingOption(3, false),
		views.SortFilter("price_asc"),
		views.SortFilter("relevance"),
		views.ActiveFilters(sampleFilters),
		views.ActiveFilters(types.FilterState{}), // No active filters
		views.FilterTag("Electronics", "category", "Electronics"),
		views.FilterTag("$25 - $100", "price", "price"),
		views.MobileFilterToggle(),

		// Catalog page components
		views.CatalogPage(sampleProducts, sampleFilters, 42, categoryCounts, cartItems, favorites, 3),
		views.CatalogContent(sampleProducts, sampleFilters, 42, cartItems, favorites),
		views.ViewToggle("grid"),
		views.ViewToggle("list"),
		views.QuickSort("price_asc"),
		views.QuickSort("relevance"),
		views.FeaturedProductsBanner([]types.Product{sampleProduct}, cartItems, favorites),
		views.FeaturedProductCard(sampleProduct, false, false),
		views.FeaturedProductCard(sampleProduct, true, true),
		views.NoProductsFound(sampleFilters),
		views.LoadMoreButton(sampleFilters, 20),

		// Product image gallery
		views.ProductImageGallery([]string{"/img1.jpg", "/img2.jpg", "/img3.jpg"}),
		views.ProductImageGallery([]string{}), // No images
		views.ProductImageGallery([]string{"/single-img.jpg"}), // Single image

		// Quick action buttons
		views.QuickActionButton("Quick View", "eye", "/products/sample-001"),
		views.QuickActionButton("Add to Wishlist", "heart", "#"),
		views.QuickActionButton("Remove from Wishlist", "heart-filled", "#"),
	); err != nil {
		log.Fatalf("❌ Code generation failed: %v", err)
	}

	log.Println("✅ Successfully generated optimized CSS classes!")
	log.Println("📁 Generated files:")
	log.Println("   - classes/classes.go (Go class constants)")
	log.Println("   - input.css (CSS with @apply rules)")
	log.Println("   - classes/classes.html (HTML for Tailwind purging)")
	log.Println("")
	log.Println("🚀 Next steps:")
	log.Println("   1. Run: tailwindcss -i input.css -o _static/dist/styles.css --watch")
	log.Println("   2. Run: templ generate")
	log.Println("   3. Run: go run main.go")
	log.Println("   4. Open: http://localhost:8080")
}

// Helper function to create pointer to bool
func ptr[T any](v T) *T {
	return &v
}