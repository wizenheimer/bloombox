package filter

import (
	"fmt"
	"sync"

	cuckoo "github.com/linvon/cuckoo-filter"
)

// CuckooFilter implements Filter using cuckoo filter
type CuckooFilter struct {
	filter   *cuckoo.Filter
	size     int
	capacity int
	mu       sync.RWMutex
}

// NewCuckooFilter creates a new cuckoo filter
func NewCuckooFilter(capacity int, falsePositiveRate float64) *CuckooFilter {
	// Calculate optimal parameters for cuckoo filter
	bucketSize := uint(4) // Standard bucket size

	// Adjust capacity based on false positive rate
	adjustedCapacity := int(float64(capacity) * 1.2) // Add 20% buffer

	filter := cuckoo.NewFilter(uint(adjustedCapacity), bucketSize, uint(capacity), cuckoo.TableTypeSingle)

	return &CuckooFilter{
		filter:   filter,
		size:     0,
		capacity: adjustedCapacity,
	}
}

// Add adds an item to the filter
func (cf *CuckooFilter) Add(item string) error {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	if cf.filter.Add([]byte(item)) {
		cf.size++
		return nil
	}
	return fmt.Errorf("failed to add item to cuckoo filter: %s", item)
}

// Contains checks if an item exists in the filter
func (cf *CuckooFilter) Contains(item string) bool {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	return cf.filter.Contain([]byte(item))
}

// Size returns the number of items in the filter
func (cf *CuckooFilter) Size() int {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	return cf.size
}

// Delete removes an item from the filter (if supported)
func (cf *CuckooFilter) Delete(item string) bool {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	if cf.filter.Delete([]byte(item)) {
		cf.size--
		return true
	}
	return false
}

// Clear removes all items from the filter
// Note: Cuckoo filters don't have a direct clear method, so we recreate the filter
func (cf *CuckooFilter) Clear() {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	// Recreate the filter with the same parameters
	bucketSize := uint(4)
	cf.filter = cuckoo.NewFilter(uint(cf.capacity), bucketSize, uint(cf.capacity), cuckoo.TableTypeSingle)
	cf.size = 0
}

// GetAll returns all items in the filter
// Note: Cuckoo filters don't support enumeration, so this returns an empty slice
func (cf *CuckooFilter) GetAll() []string {
	// Cuckoo filters don't support enumeration of stored items
	// This is a limitation of the data structure
	return []string{}
}
