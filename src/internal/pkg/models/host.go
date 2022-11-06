package models

const (
	HostTypeSourceCodeRepository = "repository"
)

type Host struct {
	Id                 string
	Name               string
	BaseUrl            string
	Type               string
	SubType            string
	AuthenticationType string
	ClientSecretName   string
}
