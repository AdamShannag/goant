package walker_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/AdamShannag/goant/pkg/walker"
)

func TestFileWalker_Walk(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{
		"main.go",
		"helper.go",
		"readme.txt",
		"nested/test.go",
		"nested/ignore.md",
	}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte("package main"), 0o644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	fileWalker := walker.NewFileWalker([]string{})

	t.Run("only .go files are visited", func(t *testing.T) {
		var visited []string
		err := fileWalker.Walk(tmpDir, func(path string) error {
			visited = append(visited, filepath.Base(path))
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := map[string]bool{
			"main.go":   true,
			"helper.go": true,
			"test.go":   true,
		}
		if len(visited) != len(expected) {
			t.Fatalf("expected %d files, got %d", len(expected), len(visited))
		}
		for _, f := range visited {
			if !expected[f] {
				t.Errorf("unexpected file visited: %s", f)
			}
		}
	})

	t.Run("callback error is returned", func(t *testing.T) {
		myErr := os.ErrPermission
		err := fileWalker.Walk(tmpDir, func(path string) error {
			return myErr
		})
		if err != nil {
			if !errors.Is(err, myErr) {
				t.Errorf("expected error %v, got %v", myErr, err)
			}
		} else {
			t.Error("expected error but got nil")
		}
	})

	t.Run("ignores non-existent root", func(t *testing.T) {
		err := fileWalker.Walk(filepath.Join(tmpDir, "doesnotexist"), func(path string) error {
			t.Error("callback should not be called")
			return nil
		})
		if err != nil {
			if !os.IsNotExist(err) {
				t.Errorf("expected not exists error, got %v", err)
			}
		}
	})
}
