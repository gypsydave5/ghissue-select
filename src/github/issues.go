package github

import "github.com/gypsydave5/ghissue-select/src"

type IssuesRepository struct{}

func (r IssuesRepository) CreateIssue() src.Issue {
	return 0
}

func (r IssuesRepository) GetIssues() []src.Issue {
	return []src.Issue{}
}

func NewIssuesRepository() *IssuesRepository {
	return &IssuesRepository{}
}
