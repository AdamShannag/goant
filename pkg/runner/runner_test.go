package runner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/AdamShannag/goant/pkg/runner"
)

func TestRun_Success(t *testing.T) {
	called := false
	ex := func(name string, args ...string) func() error {
		called = true
		if name != "echo" || len(args) != 1 || args[0] != "hello" {
			t.Fatalf("unexpected command: %s %v", name, args)
		}
		return func() error { return nil }
	}

	r := &runner.CommandRunnerImpl{Executor: ex}

	err := r.Run("echo @msg", map[string]string{"msg": "hello"}, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("executor was not called")
	}
}

func TestRun_ErrorFromCommand(t *testing.T) {
	ex := func(name string, args ...string) func() error {
		return func() error { return fmt.Errorf("command failed") }
	}

	r := &runner.CommandRunnerImpl{Executor: ex}
	err := r.Run("ls @opt", map[string]string{"opt": "-z"}, false)
	if err == nil || err.Error() != "command failed" {
		t.Fatalf("expected command failed error, got %v", err)
	}
}

func TestRun_EmptyTemplate(t *testing.T) {
	r := runner.NewCommandRunner()
	err := r.Run("", nil, false)
	if err == nil || !strings.Contains(err.Error(), "template cannot be empty") {
		t.Fatalf("expected empty template error, got %v", err)
	}
}

func TestRun_InvalidAfterReplacement(t *testing.T) {
	r := runner.NewCommandRunner()
	err := r.Run("   ", nil, false)
	if err == nil || !strings.Contains(err.Error(), "invalid command") {
		t.Fatalf("expected invalid command error, got %v", err)
	}
}

func TestRun_DryRun(t *testing.T) {
	called := false
	ex := func(name string, args ...string) func() error {
		called = true
		return func() error { return nil }
	}

	r := &runner.CommandRunnerImpl{Executor: ex}
	err := r.Run("echo @a @b", map[string]string{"a": "one", "b": "two"}, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if called {
		t.Fatal("executor should not be called in dry-run mode")
	}
}

func TestRun_MultiplePlaceholders(t *testing.T) {
	executedCmd := ""
	ex := func(name string, args ...string) func() error {
		executedCmd = name + " " + strings.Join(args, " ")
		return func() error { return nil }
	}

	r := &runner.CommandRunnerImpl{Executor: ex}
	err := r.Run("echo @first @second", map[string]string{"second": "two", "first": "one"}, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "echo one two"
	if executedCmd != expected {
		t.Fatalf("expected command '%s', got '%s'", expected, executedCmd)
	}
}
