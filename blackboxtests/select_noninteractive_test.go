package blackboxtests

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"os/exec"
	"testing"
)

func Test_NonInteractiveSelectHook_WhenSomeoneIs_WorkingAlone(t *testing.T) {
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

func Test_NonInteractiveSelectHook_WhenSomeoneIs_WorkingAlone_AndThereIsNoPairsFile(t *testing.T) {
	t.Cleanup(cleanup)

	var (
		commitMessage = "feat-376 Did some work"
	)
	givenThereIsACommitMessageFile(t, commitMessage)
	givenThereIsNotAPairsFile()

	_, err := runNonInteractiveSelectHook(t)
	assert.NoError(t, err)

	expectedMessage := commitMessage
	assertCommitMessageFileHasContents(t, expectedMessage)
	assertNoIssueFile(t)
}

func Test_NonInteractiveSelectHook_WhenSomeoneIs_Pairing_WithASinglePerson(t *testing.T) {
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

func runNonInteractiveSelectHook(t *testing.T) (string, error) {
	t.Helper()
	cmd := exec.Command(
		"go", "run", "../cmd/select/...",
		fmt.Sprintf("--commitFile=%s", commitFilePath),
		fmt.Sprintf("--pairsFile=%s", issueFilePath),
	)

	b, err := cmd.CombinedOutput()
	t.Log("CLI output:\n", string(b))
	return stripansi.Strip(string(b)), err
}
