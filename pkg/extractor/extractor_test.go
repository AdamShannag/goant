package extractor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AdamShannag/goant/pkg/extractor"
)

func TestExtract_WithAnnotation(t *testing.T) {
	src := `package test

// gomock:package=users out=.gen/mocks/repo.go
type Repository interface {
	Find(id string) (string, error)
}`

	file := writeTempFile(t, src)
	e := extractor.NewTypeExtractor()

	result, err := e.Extract(file, "gomock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	ann := result[0]
	if ann.TypeName != "Repository" {
		t.Errorf("expected type Repository, got %s", ann.TypeName)
	}
	if ann.FilePath != file {
		t.Errorf("expected path %s, got %s", file, ann.FilePath)
	}

	if ann.Params["package"] != "users" {
		t.Errorf("expected package=users, got %s", ann.Params["package"])
	}
	if ann.Params["out"] != ".gen/mocks/repo.go" {
		t.Errorf("expected out=.gen/mocks/repo.go, got %s", ann.Params["out"])
	}
	if ann.Params["type"] != "Repository" {
		t.Errorf("expected type=Repository, got %s", ann.Params["type"])
	}
	if ann.Params["path"] != file {
		t.Errorf("expected path=%s, got %s", file, ann.Params["path"])
	}
}

func TestExtract_NoAnnotation(t *testing.T) {
	src := `package test

type Service struct {
	Name string
}`

	file := writeTempFile(t, src)
	e := extractor.NewTypeExtractor()

	result, err := e.Extract(file, "gomock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected no results, got %d", len(result))
	}
}

func TestExtract_DifferentKeyword(t *testing.T) {
	src := `package test

// custom:foo=bar
type Foo struct{}`
	file := writeTempFile(t, src)
	e := extractor.NewTypeExtractor()

	result, err := e.Extract(file, "custom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	ann := result[0]
	if ann.Params["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %s", ann.Params["foo"])
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "example.go")

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return tmpFile
}
