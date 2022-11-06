package models

type Repository struct {
	Id           string
	Organization Organization
	Name         string
	Url          string
	RawData      interface{}
}
