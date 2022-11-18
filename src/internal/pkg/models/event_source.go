package models

type EventSource struct {
	Id      string
	Name    string
	Type    string
	Active  bool
	Secrets []*AuthenticationSecret
}

type AuthenticationSecret struct {
	SecretName string
	Active     bool
}
