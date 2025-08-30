package app_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/AdamShannag/goant/pkg/app"
	"github.com/AdamShannag/goant/pkg/extractor"
)

func TestApp_Run(t *testing.T) {
	paths := []string{"file1.go", "file2.go"}
	anns := map[string][]extractor.TypeAnnotation{
		"file1.go": {
			{TypeName: "T1", FilePath: "file1.go", Params: map[string]string{"p": "v1"}},
		},
		"file2.go": {
			{TypeName: "T2", FilePath: "file2.go", Params: map[string]string{"p": "v2"}},
		},
	}

	walker := &fakeWalker{paths: paths}
	extr := &fakeExtractor{annotations: anns}
	runner := &fakeRunner{}

	options := app.RunOptions{
		Root:    ".",
		Keyword: "gomock",
		Cmd:     "echo @p",
		DryRun:  true,
	}

	a := app.New(options, walker, extr, runner)

	err := a.Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(runner.called) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(runner.called))
	}

	expected := []string{"echo v1", "echo v2"}
	for i, cmd := range runner.called {
		if cmd != expected[i] {
			t.Errorf("expected command '%s', got '%s'", expected[i], cmd)
		}
	}
}

func TestApp_Run_WithRunnerError(t *testing.T) {
	paths := []string{"file1.go"}
	anns := map[string][]extractor.TypeAnnotation{
		"file1.go": {
			{TypeName: "T1", FilePath: "file1.go", Params: map[string]string{"p": "v1"}},
		},
	}

	walker := &fakeWalker{paths: paths}
	extr := &fakeExtractor{annotations: anns}
	runner := &fakeRunner{err: errors.New("fail")}

	options := app.RunOptions{
		Root:    ".",
		Keyword: "gomock",
		Cmd:     "echo @p",
		DryRun:  false,
	}

	a := app.New(options, walker, extr, runner)

	err := a.Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(runner.called) != 1 {
		t.Fatalf("expected 1 command, got %d", len(runner.called))
	}
	if runner.called[0] != "echo v1" {
		t.Errorf("expected command 'echo v1', got '%s'", runner.called[0])
	}
}

type fakeWalker struct {
	paths []string
}

func (f *fakeWalker) Walk(_ string, fn func(path string) error) error {
	for _, p := range f.paths {
		if err := fn(p); err != nil {
			return err
		}
	}
	return nil
}

type fakeExtractor struct {
	annotations map[string][]extractor.TypeAnnotation
}

func (f *fakeExtractor) Extract(path string, _ string) ([]extractor.TypeAnnotation, error) {
	if anns, ok := f.annotations[path]; ok {
		return anns, nil
	}
	return nil, nil
}

type fakeRunner struct {
	called []string
	err    error
}

func (f *fakeRunner) Run(template string, args map[string]string, _, _ bool) error {
	cmd := template
	for k, v := range args {
		cmd = strings.ReplaceAll(cmd, "@"+k, v)
	}
	f.called = append(f.called, cmd)
	return f.err
}
