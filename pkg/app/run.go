package app

import (
	"fmt"

	"github.com/AdamShannag/goant/pkg/extractor"
	"github.com/AdamShannag/goant/pkg/runner"
	"github.com/AdamShannag/goant/pkg/walker"
)

type RunOptions struct {
	Root    string
	Keyword string
	Cmd     string
	DryRun  bool
}

// gomock:package=users out=./gen/tests/repo.go
type App interface {
	Run() error
}

type app struct {
	options       RunOptions
	fileWalker    walker.FileWalker
	typeExtractor extractor.TypeExtractor
	commandRunner runner.CommandRunner
}

func New(
	options RunOptions,
	fileWalker walker.FileWalker,
	typeExtractor extractor.TypeExtractor,
	commandRunner runner.CommandRunner,
) App {
	return &app{
		options:       options,
		fileWalker:    fileWalker,
		typeExtractor: typeExtractor,
		commandRunner: commandRunner,
	}
}

func (a *app) Run() error {
	return a.fileWalker.Walk(a.options.Root, func(path string) error {
		annotations, err := a.typeExtractor.Extract(path, a.options.Keyword)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", path, err)
			return nil
		}

		for _, ann := range annotations {
			if err = a.commandRunner.Run(a.options.Cmd, ann.Params, a.options.DryRun); err != nil {
				fmt.Printf("Error running command for %s: %v\n", ann.TypeName, err)
			}
		}
		return nil
	})
}
