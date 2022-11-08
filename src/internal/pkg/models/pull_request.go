package models

import "time"

type PullRequest struct {
	Id         string
	Repository *Repository
	Url        string
	Title      string
	CreatedAt  time.Time
	Status     string
}
