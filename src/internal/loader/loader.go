// Package loader provides functionality for loading email lists and domains from files.
//
// The package is organized into the following files:
//   - file_loader.go: Basic file loading implementation
//   - batch_loader.go: Parallel loading capabilities
//   - file_watcher.go: File change monitoring
//   - utils.go: Utility functions for file validation and statistics
//
// Usage:
//
//	loader := NewFileLoader() // Create a new file loader
//	items, err := loader.LoadFromFile("emails.txt")
//
//	batchLoader := NewBatchLoader(100, 4) // 100 items per batch, 4 workers
//	items, err := batchLoader.LoadFromFilesParallel([]string{"file1.txt", "file2.txt"})
//
//	watcher := NewFileWatcher("config.txt")
//	if watcher.HasChanged() {
//	    items, changed, err := watcher.LoadIfChanged()
//	}
package loader

// Loader interface for loading email lists
type Loader interface {
	LoadFromFile(filename string) ([]string, error)
	LoadFromFiles(filenames []string) ([]string, error)
	LoadFromString(content string) ([]string, error)
}
