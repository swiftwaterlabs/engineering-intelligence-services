package models

import "time"

type Repository struct {
	Id                  string
	Organization        Organization
	Name                string
	Url                 string
	Type                string
	DefaultBranch       string
	ContentsLastUpdated time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Visibility          string
	IsForkedRepository  bool
	ForksCount          int
	RawData             interface{}
}
