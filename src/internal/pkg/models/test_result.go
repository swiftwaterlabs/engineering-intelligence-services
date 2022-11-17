package models

type TestResult struct {
	Id         string
	Type       string
	Host       string
	HostType   string
	Project    string
	AnalyzedAt string
	Metrics    map[string]string
	RawData    interface{}
}
