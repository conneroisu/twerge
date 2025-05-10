package twerge

import "sync"

// DebugHandler is a [Handler] that can be used to debug tailwind classes.
//
// It is not meant to be used in production.
//
// It will return the same class name for the same input.
type DebugHandler struct {
	mu    sync.RWMutex
	cache map[string]CacheValue
}

// NewDebugHandler creates a new DebugHandler.
func NewDebugHandler() *DebugHandler {
	return &DebugHandler{
		cache: make(map[string]CacheValue),
	}
}

// :GoImpl d *DebugHandler twerge.Handler

// It returns a short unique CSS class name from the merged classes.
func (d *DebugHandler) It(s string) string { return s }

// Cache returns the cache of the [Generator].
func (d *DebugHandler) Cache() map[string]CacheValue {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.cache
}

// SetCache sets the cache of the [Generator].
func (d *DebugHandler) SetCache(newC map[string]CacheValue) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache = newC
}
