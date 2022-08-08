package src_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/gypsydave5/ghissue-select/src"
	"testing"
)

func TestAddingGitIssueToPlainMessage(t *testing.T) {
	commitMessage := "Hello world :D"
	issue := 123

	expectedMessage := fmt.Sprintf("Hello world :D\n\n#%d", issue)

	preparedMessage := src.PrepareCommitMessage(commitMessage, issue)
	assert.Equal(t, expectedMessage, preparedMessage)
}

//
func TestDoesNotAddGitIssueThatAlreadyExists(t *testing.T) {
	commitMessage := "Hello world :D\n#123"
	issue := 123

	expectedMessage := commitMessage + "\n"

	preparedMessage := src.PrepareCommitMessage(commitMessage, issue)
	assert.Equal(t, expectedMessage, preparedMessage)
}

//
func TestAddingCoAuthorsToTemplatedMessage(t *testing.T) {
	inputMessage := "Hello world :D" + src.COMMIT_SEPARATOR + "\nother stuff"
	issue := 123

	expectedMessage := fmt.Sprintf("Hello world :D\n\n#%d%s\nother stuff", issue, src.COMMIT_SEPARATOR)

	actualMessage := src.PrepareCommitMessage(inputMessage, issue)
	assert.Equal(t, expectedMessage, actualMessage)
}
