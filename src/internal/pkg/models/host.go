package models

const (
	HostTypeSourceCodeRepository     = "source code host"
	HostTypeAutomatedTestingPlatform = "automated tests"
)

type Host struct {
	Id                 string
	Name               string
	BaseUrl            string
	Type               string
	SubType            string
	AuthenticationType string
	ClientId           string
	ClientSecretName   string
}
