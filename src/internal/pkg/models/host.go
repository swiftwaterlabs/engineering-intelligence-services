package models

const (
	HostTypeSourceCodeRepository = "source code host"
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
