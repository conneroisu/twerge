package data

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/conneroisu/twerge/examples/ecommerce-catalog/types"
)

type Store struct {
	products   []types.Product
	carts      map[string]types.Cart
	favorites  map[string]map[string]bool
	mutex      sync.RWMutex
}

func NewStore() *Store {
	store := &Store{
		products:  generateSampleProducts(),
		carts:     make(map[string]types.Cart),
		favorites: make(map[string]map[string]bool),
	}
	return store
}

func (s *Store) GetProducts(filters types.FilterState) []types.Product {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	products := make([]types.Product, 0)
	
	for _, product := range s.products {
		if s.matchesFilters(product, filters) {
			products = append(products, product)
		}
	}

	// Sort products
	s.sortProducts(products, filters.SortBy)

	return products
}

func (s *Store) GetProductByID(id string) (types.Product, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, product := range s.products {
		if product.ID == id {
			return product, nil
		}
	}
	return types.Product{}, fmt.Errorf("product not found")
}

func (s *Store) GetCategoryCounts() map[string]int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	counts := make(map[string]int)
	for _, product := range s.products {
		counts[product.Category]++
	}
	return counts
}

func (s *Store) GetCart(sessionID string) types.Cart {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if cart, exists := s.carts[sessionID]; exists {
		return cart
	}
	return types.Cart{
		ID:       sessionID,
		Items:    []types.CartItem{},
		Subtotal: 0,
		Tax:      0,
		Shipping: 0,
		Total:    0,
	}
}

func (s *Store) GetCartItemMap(sessionID string) map[string]bool {
	cart := s.GetCart(sessionID)
	itemMap := make(map[string]bool)
	for _, item := range cart.Items {
		itemMap[item.Product.ID] = true
	}
	return itemMap
}

func (s *Store) GetCartCount(sessionID string) int {
	cart := s.GetCart(sessionID)
	count := 0
	for _, item := range cart.Items {
		count += item.Quantity
	}
	return count
}

func (s *Store) AddToCart(sessionID, productID string, quantity int, variants []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	product, err := s.getProductByIDUnsafe(productID)
	if err != nil {
		return err
	}

	cart := s.getCartUnsafe(sessionID)
	
	// Check if item already exists
	for i, item := range cart.Items {
		if item.Product.ID == productID {
			cart.Items[i].Quantity += quantity
			s.updateCartTotals(&cart)
			s.carts[sessionID] = cart
			return nil
		}
	}

	// Add new item
	item := types.CartItem{
		ID:               fmt.Sprintf("%s-%d", productID, time.Now().Unix()),
		Product:          product,
		Quantity:         quantity,
		SelectedVariants: variants,
		Price:            product.Price,
		AddedAt:          time.Now(),
	}

	if product.SalePrice != nil {
		item.Price = *product.SalePrice
	}

	cart.Items = append(cart.Items, item)
	s.updateCartTotals(&cart)
	s.carts[sessionID] = cart

	return nil
}

func (s *Store) RemoveFromCart(sessionID, itemID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	cart := s.getCartUnsafe(sessionID)
	
	for i, item := range cart.Items {
		if item.ID == itemID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			s.updateCartTotals(&cart)
			s.carts[sessionID] = cart
			return nil
		}
	}

	return fmt.Errorf("item not found")
}

func (s *Store) UpdateCartItemQuantity(sessionID, itemID string, quantity int) (types.CartItem, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	cart := s.getCartUnsafe(sessionID)
	
	for i, item := range cart.Items {
		if item.ID == itemID {
			cart.Items[i].Quantity = quantity
			s.updateCartTotals(&cart)
			s.carts[sessionID] = cart
			return cart.Items[i], nil
		}
	}

	return types.CartItem{}, fmt.Errorf("item not found")
}

func (s *Store) GetFavorites(sessionID string) map[string]bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if favorites, exists := s.favorites[sessionID]; exists {
		return favorites
	}
	return make(map[string]bool)
}

// Helper methods

func (s *Store) getProductByIDUnsafe(id string) (types.Product, error) {
	for _, product := range s.products {
		if product.ID == id {
			return product, nil
		}
	}
	return types.Product{}, fmt.Errorf("product not found")
}

func (s *Store) getCartUnsafe(sessionID string) types.Cart {
	if cart, exists := s.carts[sessionID]; exists {
		return cart
	}
	return types.Cart{
		ID:       sessionID,
		Items:    []types.CartItem{},
		Subtotal: 0,
		Tax:      0,
		Shipping: 0,
		Total:    0,
	}
}

func (s *Store) updateCartTotals(cart *types.Cart) {
	cart.Subtotal = 0
	for _, item := range cart.Items {
		cart.Subtotal += item.Price * float64(item.Quantity)
	}
	
	cart.Tax = cart.Subtotal * 0.08 // 8% tax
	cart.Shipping = 0
	if cart.Subtotal > 0 && cart.Subtotal < 50 {
		cart.Shipping = 5.99
	}
	
	cart.Total = cart.Subtotal + cart.Tax + cart.Shipping
	cart.UpdatedAt = time.Now()
}

func (s *Store) matchesFilters(product types.Product, filters types.FilterState) bool {
	// Category filter
	if len(filters.Categories) > 0 {
		found := false
		for _, category := range filters.Categories {
			if product.Category == category {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Price filter
	price := product.Price
	if product.SalePrice != nil {
		price = *product.SalePrice
	}
	if price < filters.PriceRange[0] || price > filters.PriceRange[1] {
		return false
	}

	// Stock filter
	if filters.InStock != nil && *filters.InStock && !product.InStock {
		return false
	}

	// Sale filter
	if filters.OnSale != nil && *filters.OnSale && product.SalePrice == nil {
		return false
	}

	// Rating filter
	if len(filters.Ratings) > 0 {
		found := false
		for _, rating := range filters.Ratings {
			if product.Rating >= float64(rating) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Search filter
	if filters.SearchQuery != "" {
		query := strings.ToLower(filters.SearchQuery)
		if !strings.Contains(strings.ToLower(product.Name), query) &&
		   !strings.Contains(strings.ToLower(product.Description), query) &&
		   !strings.Contains(strings.ToLower(product.Category), query) {
			return false
		}
	}

	return true
}

func (s *Store) sortProducts(products []types.Product, sortBy string) {
	switch sortBy {
	case "price_asc":
		sort.Slice(products, func(i, j int) bool {
			priceI := products[i].Price
			if products[i].SalePrice != nil {
				priceI = *products[i].SalePrice
			}
			priceJ := products[j].Price
			if products[j].SalePrice != nil {
				priceJ = *products[j].SalePrice
			}
			return priceI < priceJ
		})
	case "price_desc":
		sort.Slice(products, func(i, j int) bool {
			priceI := products[i].Price
			if products[i].SalePrice != nil {
				priceI = *products[i].SalePrice
			}
			priceJ := products[j].Price
			if products[j].SalePrice != nil {
				priceJ = *products[j].SalePrice
			}
			return priceI > priceJ
		})
	case "rating":
		sort.Slice(products, func(i, j int) bool {
			return products[i].Rating > products[j].Rating
		})
	case "newest":
		sort.Slice(products, func(i, j int) bool {
			return products[i].CreatedAt.After(products[j].CreatedAt)
		})
	case "name_asc":
		sort.Slice(products, func(i, j int) bool {
			return products[i].Name < products[j].Name
		})
	case "name_desc":
		sort.Slice(products, func(i, j int) bool {
			return products[i].Name > products[j].Name
		})
	default: // relevance
		// Keep original order or implement relevance scoring
	}
}