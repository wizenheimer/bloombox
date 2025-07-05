package loader

import (
	"os"
)

// FileWatcher provides file watching capabilities (basic implementation)
type FileWatcher struct {
	filename     string
	lastModified int64
	loader       *FileLoader
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(filename string) *FileWatcher {
	info, _ := os.Stat(filename)
	var lastMod int64
	if info != nil {
		lastMod = info.ModTime().Unix()
	}

	return &FileWatcher{
		filename:     filename,
		lastModified: lastMod,
		loader:       NewFileLoader(),
	}
}

// HasChanged checks if the file has been modified
func (fw *FileWatcher) HasChanged() bool {
	info, err := os.Stat(fw.filename)
	if err != nil {
		return false
	}

	currentMod := info.ModTime().Unix()
	if currentMod > fw.lastModified {
		fw.lastModified = currentMod
		return true
	}

	return false
}

// LoadIfChanged loads the file only if it has been modified
func (fw *FileWatcher) LoadIfChanged() ([]string, bool, error) {
	if !fw.HasChanged() {
		return nil, false, nil
	}

	items, err := fw.loader.LoadFromFile(fw.filename)
	if err != nil {
		return nil, false, err
	}

	return items, true, nil
}
