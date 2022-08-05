package lib

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Issue = int

const COMMIT_SEPARATOR = "\n# ------------------------ >8 ------------------------"

func PrepareCommitMessage(input string, issue int) string {
	sections := strings.SplitN(input, COMMIT_SEPARATOR, 2)
	message := sections[0] + "\n"
	issueString := "#" + strconv.Itoa(issue)

	if !strings.Contains(message, issueString) {
		message += "\n" + issueString
	}

	if len(sections) > 1 {
		metadataSection := sections[1]
		message = message + COMMIT_SEPARATOR + metadataSection
	}

	return message
}

func LoadCommitMessage(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read commit message file %q - %w", filePath, err)
	}
	return string(b), nil
}

func DoesCommitContainCoAuthors(input string) bool {
	return strings.Contains(input, "Co-authored-by: ")
}
