package models

type Organization struct {
	Id          string
	Host        string
	Url         string
	Name        string
	Description string
	RawData     interface{}
}
