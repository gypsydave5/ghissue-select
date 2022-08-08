package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gypsydave5/ghissue-select/src"
	"github.com/manifoldco/promptui"
	"os"
	"strconv"
	"strings"
)

type IssuesSelector struct {
	issues  []src.Issue
	options selectOptions
}

func newIssues(options selectOptions, is []src.Issue) *IssuesSelector {
	return &IssuesSelector{issues: is, options: options}
}

func (i IssuesSelector) Get(ctx context.Context) (src.Issue, bool, error) {
	return i.getIssueInteractive()
}

func (i IssuesSelector) getIssueInteractive() (src.Issue, bool, error) {
	previousIssue, wantsToUsePreviousIssue, err := i.getPreviousIssueInteractive()
	if err != nil {
		return 0, false, err
	}

	if wantsToUsePreviousIssue {
		return previousIssue, true, nil
	}

	issue, ok, err := i.getIssueNameInteractive()
	if err != nil {
		return 0, false, err
	}
	return issue, ok, nil
}

func (i IssuesSelector) getPreviousIssueInteractive() (src.Issue, bool, error) {
	var issue src.Issue
	issueFile, err := os.ReadFile(i.options.issueFilePath)
	if err != nil {
		return 0, false, nil
	}

	if err = json.NewDecoder(bytes.NewReader(issueFile)).Decode(&issue); err != nil {
		return 0, false, fmt.Errorf("failed to decode issue file %q - %w", i.options.issueFilePath, err)
	}

	yesOrNo := []string{"Yes", "No"}
	prompt := promptui.Select{
		Label:             fmt.Sprintf("Are you still working on this GitHub issue? [#%d]", issue),
		Items:             []string{"Yes", "No"},
		StartInSearchMode: i.options.ForceSearchPrompts,
		Searcher:          newSearcher(yesOrNo),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return 0, false, fmt.Errorf("failed to figure out if you're still working on the last issue: %w", err)
	}

	return issue, result == "Yes", nil
}

func (i IssuesSelector) getIssueNameInteractive() (src.Issue, bool, error) {
	var issue src.Issue

	validate := func(input string) error {
		if input == "" {
			return nil
		}
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid issue")
		}
		return nil
	}

	issueSelection := promptui.Prompt{
		Label:    "GitHub issue (default none)",
		Validate: validate,
	}

	issueString, err := issueSelection.Run()
	if err != nil {
		return 0, false, fmt.Errorf("failed to input a valid issue - %w", err)
	}

	if issueString == "" {
		return 0, false, nil
	}

	issue, err = strconv.Atoi(issueString)
	if err != nil {
		return 0, false, fmt.Errorf("failed to input a valid issue - %w", err)
	}

	return issue, true, nil
}

func newSearcher(items []string) func(input string, index int) bool {
	return func(input string, index int) bool {
		name := strings.ToLower(items[index])
		return strings.Contains(name, strings.ToLower(input))
	}
}
