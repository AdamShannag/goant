package runner

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type CommandRunner interface {
	Run(template string, args map[string]string, dryRun bool) error
}

type CommandExecutor func(name string, args ...string) func() error

type CommandRunnerImpl struct {
	Executor CommandExecutor
}

func NewCommandRunner() CommandRunner {
	ex := func(name string, args ...string) func() error {
		return func() error {
			cmd := exec.Command(name, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}
	return &CommandRunnerImpl{Executor: ex}
}

func (r *CommandRunnerImpl) Run(template string, args map[string]string, dryRun bool) error {
	if template == "" {
		return fmt.Errorf("template cannot be empty")
	}

	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var replArgs []string
	for _, k := range keys {
		replArgs = append(replArgs, "@"+k, args[k])
	}

	cmdStr := strings.NewReplacer(replArgs...).Replace(template)
	fmt.Println("Running:", cmdStr)

	if dryRun {
		return nil
	}

	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("invalid command after replacement")
	}

	if r.Executor == nil {
		return fmt.Errorf("no executor defined")
	}

	return r.Executor(parts[0], parts[1:]...)()
}
