package blackboxtests

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/alecthomas/assert/v2"
	"github.com/gypsydave5/ghissue-select/src/lib"
	"os/exec"
	"testing"
)

func Test_NonInteractiveSelectHook_WorkingOnAnIssue(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
		issue         = 123
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsAnIssueFile(t, issue)

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)

	expectedMessage := lib.PrepareCommitMessage(commitMessage, issue)
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertIssueFileHasIssueEqualTo(t, issue)
}

func Test_NonInteractiveSelectHook_NotWorkingOnAnIssue(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsNoIssueFile()

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)

	expectedMessage := commitMessage
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertNoIssueFile(t)
}

func runNonInteractiveSelectHook(t *testing.T) (string, error) {
	t.Helper()
	cmd := exec.Command(
		"go", "run", "../cmd/select/...",
		fmt.Sprintf("--commitFile=%s", commitFilePath),
		fmt.Sprintf("--issueFile=%s", issueFilePath),
	)

	b, err := cmd.CombinedOutput()
	t.Log("CLI output:\n", string(b))
	return stripansi.Strip(string(b)), err
}
