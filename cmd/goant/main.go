package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AdamShannag/goant/pkg/app"
	"github.com/AdamShannag/goant/pkg/extractor"
	"github.com/AdamShannag/goant/pkg/runner"
	"github.com/AdamShannag/goant/pkg/walker"
)

const version = "0.2.0"

func main() {
	root := flag.String("root", ".", "Root directory for Go files")
	keyword := flag.String("keyword", "", "Comment keyword to search for, e.g. 'gomock'")
	cmdTemplate := flag.String("cmd", "", "Command template with placeholders, e.g. 'go run go.uber.org/mock/mockgen@latest -destination=@out -package=@package -source=@path @type'")
	dryRun := flag.Bool("dry", false, "Print the commands without executing them")
	showVersion := flag.Bool("version", false, "Print version and exit")
	silent := flag.Bool("s", false, "Suppress all output")

	flag.Parse()

	if *showVersion {
		fmt.Println("Goant version", version)
		os.Exit(0)
	}

	if *keyword == "" || *cmdTemplate == "" {
		fmt.Println("Please provide both -keyword and -cmd flags")
		flag.Usage()
		os.Exit(1)
	}

	options := app.RunOptions{
		Root:    *root,
		Keyword: *keyword,
		Cmd:     *cmdTemplate,
		DryRun:  *dryRun,
		Silent:  *silent,
	}

	if err := app.New(options,
		walker.NewFileWalker(),
		extractor.NewTypeExtractor(),
		runner.NewCommandRunner()).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
