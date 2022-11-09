package models

type Webhook struct {
	Id           string
	Type         string
	Source       string
	Organization *Organization
	Repository   *Repository
	Events       []string
	Target       string
	Active       bool
	Name         string
	RawData      interface{}
}
