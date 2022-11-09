package models

type RepositoryOwner struct {
	Id          string
	Type        string
	Repository  *Repository
	Pattern     string
	Owner       string
	ParentOwner string
}
