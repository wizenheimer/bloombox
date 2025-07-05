package loader

import (
	"sync"
)

// BatchLoader provides batch loading capabilities
type BatchLoader struct {
	loader     *FileLoader
	batchSize  int
	maxWorkers int
}

// NewBatchLoader creates a new batch loader
func NewBatchLoader(batchSize, maxWorkers int) *BatchLoader {
	return &BatchLoader{
		loader:     NewFileLoader(),
		batchSize:  batchSize,
		maxWorkers: maxWorkers,
	}
}

// LoadFromFilesParallel loads items from multiple files in parallel
func (bl *BatchLoader) LoadFromFilesParallel(filenames []string) ([]string, error) {
	type result struct {
		items []string
		err   error
	}

	results := make(chan result, len(filenames))
	semaphore := make(chan struct{}, bl.maxWorkers)

	var wg sync.WaitGroup

	for _, filename := range filenames {
		wg.Add(1)
		go func(fn string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			items, err := bl.loader.LoadFromFile(fn)
			results <- result{items: items, err: err}
		}(filename)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allItems []string
	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		allItems = append(allItems, res.items...)
	}

	// Remove duplicates
	uniqueItems := removeDuplicates(allItems)

	return uniqueItems, nil
}
