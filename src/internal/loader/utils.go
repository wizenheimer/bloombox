package loader

import (
	"fmt"
	"os"
)

// removeDuplicates removes duplicate items from slice
func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// ValidateFile checks if a file exists and is readable
func ValidateFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("file %s does not exist or is not accessible: %w", filename, err)
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", filename)
	}

	// Check if file is readable
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("file %s is not readable: %w", filename, err)
	}
	file.Close()

	return nil
}

// GetFileStats returns statistics about a file
func GetFileStats(filename string) (map[string]interface{}, error) {
	if err := ValidateFile(filename); err != nil {
		return nil, err
	}

	loader := NewFileLoader()
	items, err := loader.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}

	info, _ := os.Stat(filename)

	stats := map[string]interface{}{
		"filename":   filename,
		"size_bytes": info.Size(),
		"item_count": len(items),
		"modified":   info.ModTime(),
	}

	return stats, nil
}
