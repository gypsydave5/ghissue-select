package lib_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/tamj0rd2/coauthor-select/src/lib"
	"testing"
)

func TestAddingGitIssueToPlainMessage(t *testing.T) {
	commitMessage := "Hello world :D"
	issue := 123

	expectedMessage := fmt.Sprintf("Hello world :D\n\n#%d", issue)

	preparedMessage := lib.PrepareCommitMessage(commitMessage, issue)
	assert.Equal(t, expectedMessage, preparedMessage)
}

//
func TestDoesNotAddGitIssueThatAlreadyExists(t *testing.T) {
	commitMessage := "Hello world :D\n#123"
	issue := 123

	expectedMessage := commitMessage + "\n"

	preparedMessage := lib.PrepareCommitMessage(commitMessage, issue)
	assert.Equal(t, expectedMessage, preparedMessage)
}

//
func TestAddingCoAuthorsToTemplatedMessage(t *testing.T) {
	inputMessage := "Hello world :D" + lib.COMMIT_SEPARATOR + "\nother stuff"
	issue := 123

	expectedMessage := fmt.Sprintf("Hello world :D\n\n#%d%s\nother stuff", issue, lib.COMMIT_SEPARATOR)

	actualMessage := lib.PrepareCommitMessage(inputMessage, issue)
	assert.Equal(t, expectedMessage, actualMessage)
}
