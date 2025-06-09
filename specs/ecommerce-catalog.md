# E-commerce Product Catalog Use Case

## Overview
A modern e-commerce product catalog showcasing Twerge's dynamic styling capabilities with product variations, interactive filters, and cart animations. Demonstrates complex conditional styling and state-dependent visual feedback.

## Key Features
- Product grid with hover effects and quick-view modals
- Dynamic product variant selection (size, color, style)
- Interactive shopping cart with item animations
- Advanced filtering with visual feedback
- Price comparison and discount highlighting

## Twerge Benefits Demonstrated

### 1. Product Card with Dynamic States
```go
type Product struct {
    ID       string
    Name     string
    Price    float64
    SalePrice *float64
    InStock  bool
    Featured bool
    Images   []string
    Variants []ProductVariant
}

templ ProductCard(product Product, inCart bool, favorited bool) {
    <div class={twerge.It(
        "group relative bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-lg " +
        "transition-all duration-300 overflow-hidden border border-gray-200 dark:border-gray-700 " +
        "hover:border-gray-300 dark:hover:border-gray-600 " +
        twerge.If(product.Featured, "ring-2 ring-blue-500 ring-opacity-50", "") +
        twerge.If(!product.InStock, "opacity-75 grayscale", "")
    )}>
        <!-- Product Image with Overlay -->
        <div class={twerge.It("relative aspect-square overflow-hidden")}>
            <img 
                src={product.Images[0]} 
                alt={product.Name}
                class={twerge.It(
                    "w-full h-full object-cover transition-transform duration-500 " +
                    "group-hover:scale-110"
                )}
            />
            
            <!-- Quick Actions Overlay -->
            <div class={twerge.It(
                "absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-20 " +
                "transition-all duration-300 flex items-center justify-center"
            )}>
                <div class={twerge.It(
                    "flex space-x-2 transform translate-y-4 opacity-0 " +
                    "group-hover:translate-y-0 group-hover:opacity-100 " +
                    "transition-all duration-300"
                )}>
                    @QuickActionButton("Quick View", "eye")
                    @QuickActionButton(
                        twerge.If(favorited, "Remove from Wishlist", "Add to Wishlist"), 
                        twerge.If(favorited, "heart-filled", "heart")
                    )
                </div>
            </div>

            <!-- Stock Status Badge -->
            if !product.InStock {
                <div class={twerge.It(
                    "absolute top-2 left-2 bg-red-500 text-white text-xs font-medium " +
                    "px-2 py-1 rounded-full"
                )}>
                    Out of Stock
                </div>
            }

            <!-- Sale Badge -->
            if product.SalePrice != nil {
                <div class={twerge.It(
                    "absolute top-2 right-2 bg-red-500 text-white text-xs font-bold " +
                    "px-2 py-1 rounded-full animate-pulse"
                )}>
                    { fmt.Sprintf("%.0f%% OFF", (1-*product.SalePrice/product.Price)*100) }
                </div>
            }
        </div>

        <!-- Product Info -->
        <div class={twerge.It("p-4 space-y-3")}>
            <h3 class={twerge.It(
                "font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 " +
                "dark:group-hover:text-blue-400 transition-colors duration-200 " +
                "line-clamp-2"
            )}>
                { product.Name }
            </h3>

            <!-- Price Display -->
            <div class={twerge.It("flex items-center space-x-2")}>
                if product.SalePrice != nil {
                    <span class={twerge.It("text-lg font-bold text-red-600 dark:text-red-400")}>
                        ${ fmt.Sprintf("%.2f", *product.SalePrice) }
                    </span>
                    <span class={twerge.It("text-sm text-gray-500 line-through")}>
                        ${ fmt.Sprintf("%.2f", product.Price) }
                    </span>
                } else {
                    <span class={twerge.It("text-lg font-bold text-gray-900 dark:text-white")}>
                        ${ fmt.Sprintf("%.2f", product.Price) }
                    </span>
                }
            </div>

            <!-- Add to Cart Button -->
            <button 
                class={twerge.It(
                    "w-full py-2 px-4 rounded-lg font-medium transition-all duration-200 " +
                    twerge.If(product.InStock,
                        twerge.If(inCart,
                            "bg-green-100 text-green-700 border border-green-300 " +
                            "hover:bg-green-200 dark:bg-green-900 dark:text-green-300 " +
                            "dark:border-green-700 dark:hover:bg-green-800",
                            "bg-blue-600 text-white hover:bg-blue-700 " +
                            "active:bg-blue-800 transform hover:scale-105 active:scale-95"
                        ),
                        "bg-gray-300 text-gray-500 cursor-not-allowed dark:bg-gray-700 dark:text-gray-400"
                    )
                )}
                disabled={!product.InStock}
            >
                if !product.InStock {
                    Out of Stock
                } else if inCart {
                    ✓ In Cart
                } else {
                    Add to Cart
                }
            </button>
        </div>
    </div>
}
```

### 2. Product Variant Selector
```go
type ProductVariant struct {
    Type     string // "color", "size", "style"
    Value    string
    Available bool
    Price    *float64 // nil if same as base price
}

templ VariantSelector(variants []ProductVariant, selected string, variantType string) {
    <div class={twerge.It("space-y-3")}>
        <h4 class={twerge.It("text-sm font-medium text-gray-900 dark:text-white capitalize")}>
            { variantType }
        </h4>
        
        <div class={twerge.It(
            twerge.If(variantType == "color", 
                "flex flex-wrap gap-2",
                "grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-6 gap-2"
            )
        )}>
            for _, variant := range variants {
                if variant.Type == variantType {
                    @VariantOption(variant, selected == variant.Value, variantType)
                }
            }
        </div>
    </div>
}

templ VariantOption(variant ProductVariant, isSelected bool, variantType string) {
    <button class={twerge.It(
        "relative transition-all duration-200 " +
        twerge.If(variantType == "color",
            "w-8 h-8 rounded-full border-2 " +
            twerge.If(variant.Available,
                twerge.If(isSelected,
                    "border-gray-900 dark:border-white ring-2 ring-offset-2 ring-blue-500",
                    "border-gray-300 dark:border-gray-600 hover:border-gray-400"
                ),
                "border-gray-200 opacity-50 cursor-not-allowed"
            ),
            // Size/Style variants
            "px-3 py-2 text-sm font-medium rounded-lg border " +
            twerge.If(variant.Available,
                twerge.If(isSelected,
                    "bg-blue-600 text-white border-blue-600",
                    "bg-white dark:bg-gray-800 text-gray-900 dark:text-white " +
                    "border-gray-300 dark:border-gray-600 hover:border-gray-400 " +
                    "hover:bg-gray-50 dark:hover:bg-gray-700"
                ),
                "bg-gray-100 dark:bg-gray-800 text-gray-400 dark:text-gray-500 " +
                "border-gray-200 dark:border-gray-700 cursor-not-allowed line-through"
            )
        )
    )}
    disabled={!variant.Available}
    >
        if variantType == "color" {
            <div class={twerge.It(
                "w-full h-full rounded-full " +
                getColorClass(variant.Value)
            )}></div>
        } else {
            { variant.Value }
            if variant.Price != nil {
                <span class={twerge.It("block text-xs text-gray-500 dark:text-gray-400")}>
                    +${ fmt.Sprintf("%.2f", *variant.Price) }
                </span>
            }
        }
        
        if !variant.Available {
            <div class={twerge.It(
                "absolute inset-0 flex items-center justify-center " +
                "bg-white bg-opacity-80 dark:bg-gray-800 dark:bg-opacity-80"
            )}>
                <svg class={twerge.It("w-4 h-4 text-red-500")}>
                    <!-- X icon -->
                </svg>
            </div>
        }
    </button>
}
```

### 3. Shopping Cart with Animations
```go
templ CartItem(item CartItem, updating bool) {
    <div class={twerge.It(
        "flex items-center space-x-4 p-4 bg-white dark:bg-gray-800 rounded-lg " +
        "border border-gray-200 dark:border-gray-700 " +
        "transition-all duration-300 " +
        twerge.If(updating, "opacity-50 pointer-events-none", "")
    )}>
        <!-- Product Image -->
        <div class={twerge.It("relative w-16 h-16 flex-shrink-0 rounded-lg overflow-hidden")}>
            <img 
                src={item.Product.Images[0]} 
                alt={item.Product.Name}
                class={twerge.It("w-full h-full object-cover")}
            />
            if updating {
                <div class={twerge.It(
                    "absolute inset-0 bg-white bg-opacity-75 dark:bg-gray-800 " +
                    "dark:bg-opacity-75 flex items-center justify-center"
                )}>
                    <div class={twerge.It("animate-spin h-4 w-4 border-2 border-blue-500 border-t-transparent rounded-full")}></div>
                </div>
            }
        </div>

        <!-- Product Details -->
        <div class={twerge.It("flex-1 min-w-0 space-y-1")}>
            <h4 class={twerge.It("font-medium text-gray-900 dark:text-white truncate")}>
                { item.Product.Name }
            </h4>
            if len(item.SelectedVariants) > 0 {
                <p class={twerge.It("text-sm text-gray-500 dark:text-gray-400")}>
                    { strings.Join(item.SelectedVariants, ", ") }
                </p>
            }
            <p class={twerge.It("text-sm font-medium text-gray-900 dark:text-white")}>
                ${ fmt.Sprintf("%.2f", item.Price) }
            </p>
        </div>

        <!-- Quantity Controls -->
        <div class={twerge.It("flex items-center space-x-2")}>
            <button class={twerge.It(
                "w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-700 " +
                "hover:bg-gray-200 dark:hover:bg-gray-600 " +
                "flex items-center justify-center transition-colors duration-200 " +
                twerge.If(item.Quantity <= 1, "opacity-50 cursor-not-allowed", "")
            )}>
                -
            </button>
            
            <span class={twerge.It(
                "w-8 text-center font-medium text-gray-900 dark:text-white " +
                "transition-all duration-200 " +
                twerge.If(updating, "animate-pulse", "")
            )}>
                { strconv.Itoa(item.Quantity) }
            </span>
            
            <button class={twerge.It(
                "w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-700 " +
                "hover:bg-gray-200 dark:hover:bg-gray-600 " +
                "flex items-center justify-center transition-colors duration-200"
            )}>
                +
            </button>
        </div>

        <!-- Remove Button -->
        <button class={twerge.It(
            "text-red-500 hover:text-red-700 dark:text-red-400 " +
            "dark:hover:text-red-300 transition-colors duration-200"
        )}>
            🗑️
        </button>
    </div>
}
```

### 4. Advanced Filter System
```go
type FilterState struct {
    Categories []string
    PriceRange [2]float64
    InStock    *bool
    OnSale     *bool
    Ratings    []int
}

templ FilterSidebar(filters FilterState, totalProducts int) {
    <div class={twerge.It(
        "w-full lg:w-80 bg-white dark:bg-gray-800 rounded-lg shadow-sm " +
        "border border-gray-200 dark:border-gray-700 p-6 space-y-6"
    )}>
        <div class={twerge.It("flex items-center justify-between")}>
            <h3 class={twerge.It("text-lg font-semibold text-gray-900 dark:text-white")}>
                Filters
            </h3>
            <span class={twerge.It(
                "text-sm text-gray-500 dark:text-gray-400 " +
                "animate-pulse transition-opacity duration-300"
            )}>
                { strconv.Itoa(totalProducts) } products
            </span>
        </div>

        @CategoryFilter(filters.Categories)
        @PriceRangeFilter(filters.PriceRange)
        @AvailabilityFilter(filters.InStock, filters.OnSale)
        @RatingFilter(filters.Ratings)
    </div>
}

templ CategoryFilter(selected []string) {
    <div class={twerge.It("space-y-3")}>
        <h4 class={twerge.It("font-medium text-gray-900 dark:text-white")}>Category</h4>
        <div class={twerge.It("space-y-2")}>
            for _, category := range categories {
                @FilterCheckbox(
                    category,
                    slices.Contains(selected, category),
                    len(getProductsInCategory(category))
                )
            }
        </div>
    </div>
}

templ FilterCheckbox(label string, checked bool, count int) {
    <label class={twerge.It(
        "flex items-center justify-between cursor-pointer group " +
        "p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 " +
        "transition-colors duration-200"
    )}>
        <div class={twerge.It("flex items-center space-x-3")}>
            <input 
                type="checkbox" 
                checked={checked}
                class={twerge.It(
                    "w-4 h-4 text-blue-600 border-gray-300 dark:border-gray-600 " +
                    "rounded focus:ring-blue-500 dark:focus:ring-blue-600 " +
                    "dark:bg-gray-700 transition-colors duration-200"
                )}
            />
            <span class={twerge.It(
                "text-sm text-gray-700 dark:text-gray-300 " +
                "group-hover:text-gray-900 dark:group-hover:text-white " +
                "transition-colors duration-200"
            )}>
                { label }
            </span>
        </div>
        <span class={twerge.It(
            "text-xs text-gray-500 dark:text-gray-400 " +
            "bg-gray-100 dark:bg-gray-600 px-2 py-1 rounded-full"
        )}>
            { strconv.Itoa(count) }
        </span>
    </label>
}
```

## Performance Benefits

### Class Consolidation Example
```css
/* twerge:begin */
/* Complex product card states consolidated */
.tw-1 { @apply group relative bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-lg transition-all duration-300 overflow-hidden border border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600; }
.tw-2 { @apply group relative bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-lg transition-all duration-300 overflow-hidden border border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 ring-2 ring-blue-500 ring-opacity-50; }
.tw-3 { @apply relative aspect-square overflow-hidden; }
.tw-4 { @apply w-full h-full object-cover transition-transform duration-500 group-hover:scale-110; }
/* twerge:end */
```

### Dynamic State Management
- **Cart animations**: Smooth transitions without class conflicts
- **Variant selection**: Clean state-dependent styling
- **Filter feedback**: Real-time visual updates
- **Loading states**: Consistent loading indicators

## Code Generation Integration

```go
//go:build ignore

package main

import (
    "github.com/conneroisu/twerge"
    "github.com/yourproject/shop/views"
    "github.com/yourproject/shop/types"
)

func main() {
    // Sample data for code generation
    sampleProduct := types.Product{
        Name: "Sample Product",
        Price: 99.99,
        InStock: true,
        Featured: true,
        Images: []string{"/sample.jpg"},
    }

    if err := twerge.CodeGen(
        twerge.Default(),
        "styles/generated.go",
        "styles/input.css",
        "styles/classes.html",
        views.ProductCard(sampleProduct, false, false),
        views.ProductCard(sampleProduct, true, true),
        views.VariantSelector([]types.ProductVariant{}, "", "color"),
        views.CartItem(types.CartItem{}, false),
        views.CartItem(types.CartItem{}, true),
        views.FilterSidebar(types.FilterState{}, 42),
    ); err != nil {
        panic(err)
    }
}
```

## Expected Outcomes

1. **Enhanced UX**: Smooth animations and state transitions without performance overhead
2. **Maintainable Code**: Consistent styling patterns across product variations
3. **Performance**: Reduced CSS bundle size despite complex interactive states
4. **Scalability**: Easy to add new product types and filter options