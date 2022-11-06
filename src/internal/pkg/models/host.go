package models

type Host struct {
	Id                 string
	Name               string
	BaseUrl            string
	Type               string
	AuthenticationType string
	ClientSecretName   string
}
