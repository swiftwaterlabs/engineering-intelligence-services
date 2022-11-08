package models

import "time"

type PullRequest struct {
	Id           string
	Type         string
	Repository   *Repository
	TargetBranch string
	Url          string
	Title        string
	CreatedBy    string
	CreatedAt    time.Time
	ClosedAt     time.Time
	IsMerged     bool
	Status       string
	Reviews      []*PullRequestReview
	Files        []string
	RawData      interface{}
}

type PullRequestReview struct {
	Reviewer   string
	Status     string
	ReviewedAt time.Time
}
