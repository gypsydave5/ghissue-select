package main

import (
	"flag"
	"github.com/mattn/go-isatty"
	"os"
)

type selectOptions struct {
	CommitFilePath     string
	AuthorsFilePath    string
	issueFilePath      string
	ForceSearchPrompts bool
	Interactive        bool
}

func parseOptions() selectOptions {
	var (
		options = selectOptions{
			Interactive: isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
		}
	)

	flag.StringVar(&options.CommitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flag.StringVar(&options.issueFilePath, "issueFile", ".ghissue", "path to file with the last GitHub issue")
	flag.BoolVar(&options.ForceSearchPrompts, "forceSearchPrompts", false, "makes all prompts searches for ease of testing")
	flag.BoolVar(&options.Interactive, "interactive", options.Interactive, "whether you're using an interactive terminal")
	flag.Parse()

	return options
}
