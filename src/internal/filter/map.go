package filter

import (
	"sync"
)

// MapFilter implements Filter using a regular map (zero false positives)
type MapFilter struct {
	items map[string]bool
	mu    sync.RWMutex
}

// NewMapFilter creates a new map filter
func NewMapFilter() *MapFilter {
	return &MapFilter{
		items: make(map[string]bool),
	}
}

// Add adds an item to the filter
func (mf *MapFilter) Add(item string) error {
	mf.mu.Lock()
	defer mf.mu.Unlock()

	mf.items[item] = true
	return nil
}

// Contains checks if an item exists in the filter
func (mf *MapFilter) Contains(item string) bool {
	mf.mu.RLock()
	defer mf.mu.RUnlock()

	return mf.items[item]
}

// Size returns the number of items in the filter
func (mf *MapFilter) Size() int {
	mf.mu.RLock()
	defer mf.mu.RUnlock()

	return len(mf.items)
}

// Delete removes an item from the filter
func (mf *MapFilter) Delete(item string) bool {
	mf.mu.Lock()
	defer mf.mu.Unlock()

	if mf.items[item] {
		delete(mf.items, item)
		return true
	}
	return false
}

// Clear removes all items from the filter
func (mf *MapFilter) Clear() {
	mf.mu.Lock()
	defer mf.mu.Unlock()

	mf.items = make(map[string]bool)
}

// GetAll returns all items in the filter (useful for debugging)
func (mf *MapFilter) GetAll() []string {
	mf.mu.RLock()
	defer mf.mu.RUnlock()

	items := make([]string, 0, len(mf.items))
	for item := range mf.items {
		items = append(items, item)
	}
	return items
}
