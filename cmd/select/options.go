package main

import (
	"flag"
)

type selectOptions struct {
	CommitFilePath     string
	AuthorsFilePath    string
	issueFilePath      string
	ForceSearchPrompts bool
}

func parseOptions() selectOptions {
	var (
		options = selectOptions{}
	)

	flag.StringVar(&options.CommitFilePath, "commitFile", ".git/COMMIT_EDITMSG", "path to commit message file")
	flag.StringVar(&options.issueFilePath, "issueFile", ".ghissue", "path to file with the last GitHub issue")
	flag.BoolVar(&options.ForceSearchPrompts, "forceSearchPrompts", false, "makes all prompts searches for ease of testing")
	flag.Parse()

	return options
}
