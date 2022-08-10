package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/gypsydave5/ghissue-select/src"
	"io"
	"testing"
	"time"
)

func Test_SelectingIssueWhenThereAreNoIssues(t *testing.T) {
	is := NewInteractiveSelect(t)

	selector := NewIssuesSelector(selectOptions{ForceSearchPrompts: true}, []src.Issue{}, is.Stdin, is.Stdout)

	is.run([]string{})

	_, ok, err := selector.Get(context.Background())
	assert.NoError(t, err)
	assert.False(t, ok)
}

func Test_SelectingIssueWhenThereIsAnIssue(t *testing.T) {
	interactiveSelect := NewInteractiveSelect(t)

	selector := NewIssuesSelector(selectOptions{ForceSearchPrompts: true}, []src.Issue{123}, interactiveSelect.Stdin, interactiveSelect.Stdout)

	interactiveSelect.run([]string{"123"})

	issue, ok, err := selector.Get(context.Background())
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, issue, src.Issue(123))
}

type InteractiveSelect struct {
	t          *testing.T
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	Stdin      io.ReadCloser
	Stdout     io.WriteCloser
	StdoutChan chan string
}

func NewInteractiveSelect(t *testing.T) *InteractiveSelect {
	internalStdout, externalStdout := io.Pipe()
	externalStdin, internalStdin := io.Pipe()
	return &InteractiveSelect{
		t:      t,
		stdout: internalStdout,
		stdin:  internalStdin,
		Stdout: externalStdout,
		Stdin:  externalStdin,
	}
}

func (is *InteractiveSelect) run(textToSubmit []string) {
	is.t.Helper()

	go func() {
		maxIndex := len(textToSubmit) - 1
		for i, text := range textToSubmit {
			if _, err := io.WriteString(is.stdin, text+"\n"); err != nil {
				panic(fmt.Errorf("failed to write %q to stdin: %v\n", text, err))
			}

			if i < maxIndex {
				// the console thing promptui uses is apparently too slow to read inputs so quickly
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		b, err := io.ReadAll(is.stdout)
		if err != nil {
			is.t.Fatal(err)
		}
		is.t.Log("CLI output:\n", string(b))
		is.StdoutChan <- string(b)
	}()
}

type noopWriterCloser struct {
	io.Writer
}

func NopWriterCloser(writer io.Writer) io.WriteCloser {
	return &noopWriterCloser{Writer: writer}
}

func (wc noopWriterCloser) Close() error {
	return nil
}
