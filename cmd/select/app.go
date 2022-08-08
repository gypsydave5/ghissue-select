package main

import (
	"context"
	"fmt"
	"github.com/gypsydave5/ghissue-select/src"
)

type GetIssue func(ctx context.Context) (src.Issue, bool, error)

type SaveIssue func(ctx context.Context, pairs src.Issue) error

type FormatCommitMessage func(authors src.Issue) (string, error)

type SaveCommitMessage func(ctx context.Context, message string) error

type CLIApp struct {
	getIssue            GetIssue
	saveIssue           SaveIssue
	formatCommitMessage FormatCommitMessage
	saveCommitMessage   SaveCommitMessage
}

func NewCLIApp(
	getIssue GetIssue,
	saveIssue SaveIssue,
	formatCommitMessage FormatCommitMessage,
	saveCommitMessage SaveCommitMessage,
) *CLIApp {
	return &CLIApp{
		getIssue:            getIssue,
		saveIssue:           saveIssue,
		formatCommitMessage: formatCommitMessage,
		saveCommitMessage:   saveCommitMessage,
	}
}

func (c CLIApp) Run(ctx context.Context) error {
	issue, ok, err := c.getIssue(ctx)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	if !ok {
		return nil
	}

	if err := c.saveIssue(ctx, issue); err != nil {
		// it's really not the end of the world. No need to kill the program.
		return fmt.Errorf("failed to save issue: %w", err)
	}

	commitMessage, err := c.formatCommitMessage(issue)
	if err != nil {
		return fmt.Errorf("failed to format commit message: %w", err)
	}

	if err = c.saveCommitMessage(ctx, commitMessage); err != nil {
		return fmt.Errorf("failed to save commit message: %w", err)
	}

	fmt.Println("Added issue:", issue)
	return nil
}
