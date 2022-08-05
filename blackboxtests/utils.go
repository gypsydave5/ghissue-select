package blackboxtests

import (
	"encoding/json"
	"github.com/alecthomas/assert/v2"
	_ "github.com/gypsydave5/ghissue-select/src/lib"

	"os"
	"testing"
)

func TestMain(m *testing.M) {
	cleanup()
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func givenThereIsACommitMessageFile(t *testing.T, message string) {
	t.Helper()
	err := os.WriteFile(commitFilePath, []byte(message), 0666)
	assert.NoError(t, err)
}

func givenThereIsAnIssueFile(t *testing.T, issue int) {
	t.Helper()

	b, err := json.Marshal(issue)
	assert.NoError(t, err, "could not marshall issue")

	err = os.WriteFile(issueFilePath, b, 0666)
	assert.NoError(t, err, "could not write issue file")
}

func givenThereIsNoIssueFile() {
	_ = os.Remove(issueFilePath)
}

func assertIssueFileHasIssueEqualTo(t *testing.T, expectedIssueNumber int) {
	t.Helper()
	b, err := os.ReadFile(issueFilePath)
	assert.NoError(t, err, "could not read file %q", issueFilePath)

	var actualIssueNumber int
	assert.NoError(t, json.Unmarshal(b, &actualIssueNumber))
	assert.Equal(t, expectedIssueNumber, actualIssueNumber)
}

func assertNoIssueFile(t *testing.T) {
	t.Helper()
	_, err := os.Stat(issueFilePath)
	assert.Error(t, err)
}

func assertCommitMessageFileHasContents(t *testing.T, message string) {
	t.Helper()
	fileContent, err := os.ReadFile(commitFilePath)
	assert.NoError(t, err, "could not read commit file %q", commitFilePath)
	assert.Equal(t, message, string(fileContent))
}

func assertCommitMessageFileContainsContents(t *testing.T, message string) {
	t.Helper()
	fileContent, err := os.ReadFile(commitFilePath)
	assert.NoError(t, err, "could not read commit file %q", commitFilePath)
	assert.Contains(t, string(fileContent), message)
}

const (
	commitFilePath = "test_commit_file"
	issueFilePath  = "test_issue"
)

func cleanup() {
	_ = os.Remove(commitFilePath)
	_ = os.Remove(issueFilePath)
}
