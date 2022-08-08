package github

import "github.com/gypsydave5/ghissue-select/src"

type IssuesRepository struct{}

func (r IssuesRepository) CreateIssue() src.Issue {
	return 0
}

func NewIssuesRepository() *IssuesRepository {
	return &IssuesRepository{}
}
