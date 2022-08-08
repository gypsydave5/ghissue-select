package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gypsydave5/ghissue-select/src"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	options := parseOptions()

	i := newIssues(options)

	cliApp := NewCLIApp(
		i.Get,
		func(ctx context.Context, issue src.Issue) error {
			b, err := json.Marshal(issue)
			if err != nil {
				return fmt.Errorf("failed to marshal issue - %w", err)
			}

			if err := os.WriteFile(options.issueFilePath, b, 0644); err != nil {
				return fmt.Errorf("failed to write issue file - %w", err)
			}
			return nil
		},
		func(issue src.Issue) (string, error) {
			file, err := os.ReadFile(options.CommitFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to read commit message file: %w", err)
			}
			return src.PrepareCommitMessage(string(file), issue), nil
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
