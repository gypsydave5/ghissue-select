package blackboxtests

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/alecthomas/assert/v2"
	"github.com/gypsydave5/ghissue-select/src"
	"io"
	"os/exec"
	"testing"
	"time"
)

func Test_InteractiveSelectHook_NotWorkingOnAnIssue(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
	)

	givenThereIsACommitMessageFile(t, commitMessage)

	_, err := runInteractiveSelectHook(t, []string{""})
	assert.NoError(t, err)

	expectedMessage := commitMessage
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertNoIssueFile(t)
}

func Test_InteractiveSelectHook_StartingAnIssue(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		issue         = 123
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsNoIssueFile()

	_, err := runInteractiveSelectHook(t, []string{"123"})
	assert.NoError(t, err)

	expectedMessage := src.PrepareCommitMessage(commitMessage, issue)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertIssueFileHasIssueEqualTo(t, issue)
}

func Test_InteractiveSelectHook_StayingOnTheSameIssue(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		issue         = 123
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnIssueFile(t, issue)

	_, err := runInteractiveSelectHook(t, []string{"Yes"})
	assert.NoError(t, err)

	expectedMessage := src.PrepareCommitMessage(commitMessage, issue)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertIssueFileHasIssueEqualTo(t, issue)
}

func Test_InteractiveSelectHook_ChangingTheIssue(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		previousIssue = 123
		expectedIssue = 456
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnIssueFile(t, previousIssue)

	_, err := runInteractiveSelectHook(t, []string{"No", "456", "No one else"})
	assert.NoError(t, err)

	expectedMessage := src.PrepareCommitMessage(commitMessage, expectedIssue)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertIssueFileHasIssueEqualTo(t, expectedIssue)
}

func runInteractiveSelectHook(t *testing.T, textToSubmit []string) (string, error) {
	t.Helper()
	cmd := exec.Command(
		"go", "run", "../cmd/select/...",
		fmt.Sprintf("--commitFile=%s", commitFilePath),
		fmt.Sprintf("--issueFile=%s", issueFilePath),
		fmt.Sprintf("--forceSearchPrompts=%t", true),
	)

	cmdStdin, err := cmd.StdinPipe()
	assert.NoError(t, err)

	var rerr error

	go func() {
		defer func() {
			_ = cmdStdin.Close()
		}()

		maxIndex := len(textToSubmit) - 1
		for i, text := range textToSubmit {
			if _, err := io.WriteString(cmdStdin, text+"\n"); err != nil {
				panic(fmt.Errorf("failed to write %q to stdin: %v\n", text, err))
			}

			if i < maxIndex {
				// the console thing promptui uses is apparently too slow to read inputs so quickly
				time.Sleep(time.Second)
			}
		}
	}()

	b, err := cmd.CombinedOutput()
	t.Log("CLI output:\n", string(b))
	return stripansi.Strip(string(b)), rerr
}
