package models

type Organization struct {
	Id       string
	Type     string
	Host     string
	HostType string
	Url      string
	Name     string
	RawData  interface{}
}
