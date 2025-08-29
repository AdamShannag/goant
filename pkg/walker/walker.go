package walker

import (
	"os"
	"path/filepath"
	"strings"
)

type FileWalker interface {
	Walk(root string, fn func(path string) error) error
}

type fileWalker struct{}

func NewFileWalker() FileWalker {
	return &fileWalker{}
}

func (d *fileWalker) Walk(root string, fn func(path string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return nil
		}
		return fn(path)
	})
}
