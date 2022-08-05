package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gypsydave5/ghissue-select/src/lib"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	options := parseOptions()

	cliApp := NewCLIApp(
		func(ctx context.Context) (lib.Issue, bool, error) {

			if !options.Interactive {
				return getIssueNonInteractive(options.issueFilePath)
			}

			return getIssueInteractive(options)
		},
		func(ctx context.Context, issue lib.Issue) error {
			b, err := json.Marshal(issue)
			if err != nil {
				return fmt.Errorf("failed to marshal issue - %w", err)
			}

			if err := os.WriteFile(options.issueFilePath, b, 0644); err != nil {
				return fmt.Errorf("failed to write issue file - %w", err)
			}
			return nil
		},
		func(issue lib.Issue) (string, error) {
			file, err := os.ReadFile(options.CommitFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to read commit message file: %w", err)
			}
			return lib.PrepareCommitMessage(string(file), issue), nil
		},
		func(ctx context.Context, message string) error {
			if err := os.WriteFile(options.CommitFilePath, []byte(message), 0644); err != nil {
				return fmt.Errorf("failed to write commit message file to %q: %w", options.CommitFilePath, err)
			}
			return nil
		},
	)

	if err := cliApp.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func getIssueInteractive(options selectOptions) (lib.Issue, bool, error) {
	previousIssue, wantsToUsePreviousIssue, err := getPreviousIssueInteractive(options)
	if err != nil {
		return 0, false, err
	}

	if wantsToUsePreviousIssue {
		return previousIssue, true, nil
	}

	issue, ok, err := getIssueNameInteractive()
	if err != nil {
		return 0, false, err
	}
	return issue, ok, nil
}

func getPreviousIssueInteractive(options selectOptions) (lib.Issue, bool, error) {
	var issue lib.Issue
	issueFile, err := os.ReadFile(options.issueFilePath)
	if err != nil {
		return 0, false, nil
	}

	if err = json.NewDecoder(bytes.NewReader(issueFile)).Decode(&issue); err != nil {
		return 0, false, fmt.Errorf("failed to decode issue file %q - %w", options.issueFilePath, err)
	}

	yesOrNo := []string{"Yes", "No"}
	prompt := promptui.Select{
		Label:             fmt.Sprintf("Are you still working on this GitHub issue? [#%d]", issue),
		Items:             []string{"Yes", "No"},
		StartInSearchMode: options.ForceSearchPrompts,
		Searcher:          newSearcher(yesOrNo),
	}
	_, result, err := prompt.Run()
	if err != nil {
		return 0, false, fmt.Errorf("failed to figure out if you're still working on the last issue: %w", err)
	}

	return issue, result == "Yes", nil
}

func getIssueNameInteractive() (lib.Issue, bool, error) {
	var issue lib.Issue

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

func getIssueNonInteractive(issueFilePath string) (lib.Issue, bool, error) {
	issue, err := readJSON[int](issueFilePath)
	if err != nil {
		return 0, false, nil
	}
	return issue, true, err
}

func readJSON[T any](filePath string) (T, error) {
	var result T
	b, err := os.ReadFile(filePath)
	if err != nil {
		return result, fmt.Errorf("failed to read file %q - %w", filePath, err)
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal file %q - %w", filePath, err)
	}

	return result, nil
}
