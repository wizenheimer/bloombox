package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// FileLoader implements email list loading from files
type FileLoader struct {
	mu sync.RWMutex
}

// NewFileLoader creates a new file loader
func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// LoadFromFile loads domains/emails from a text file
func (fl *FileLoader) LoadFromFile(filename string) ([]string, error) {
	fl.mu.RLock()
	defer fl.mu.RUnlock()

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var items []string
	scanner := bufio.NewScanner(file)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove inline comments
		if commentIndex := strings.Index(line, "#"); commentIndex != -1 {
			line = strings.TrimSpace(line[:commentIndex])
		}

		// Clean and validate the item
		item := strings.ToLower(line)
		if item != "" {
			items = append(items, item)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}

	return items, nil
}

// LoadFromFiles loads items from multiple files
func (fl *FileLoader) LoadFromFiles(filenames []string) ([]string, error) {
	var allItems []string

	for _, filename := range filenames {
		items, err := fl.LoadFromFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to load from file %s: %w", filename, err)
		}
		allItems = append(allItems, items...)
	}

	// Remove duplicates
	uniqueItems := removeDuplicates(allItems)

	return uniqueItems, nil
}

// LoadFromString loads items from a string content
func (fl *FileLoader) LoadFromString(content string) ([]string, error) {
	var items []string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove inline comments
		if commentIndex := strings.Index(line, "#"); commentIndex != -1 {
			line = strings.TrimSpace(line[:commentIndex])
		}

		// Clean the item
		item := strings.ToLower(line)
		if item != "" {
			items = append(items, item)
		}
	}

	return items, nil
}
