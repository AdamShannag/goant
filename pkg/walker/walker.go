package walker

import (
	"os"
	"path/filepath"
	"strings"
)

type FileWalker interface {
	Walk(root string, fn func(path string) error) error
}

type fileWalker struct {
	skipPaths []string
}

func NewFileWalker(skip []string) FileWalker {
	return &fileWalker{skipPaths: skip}
}

func (d *fileWalker) Walk(root string, fn func(path string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, skip := range d.skipPaths {
			if info.IsDir() && info.Name() == skip {
				return filepath.SkipDir
			}
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			return fn(path)
		}

		return nil
	})
}
