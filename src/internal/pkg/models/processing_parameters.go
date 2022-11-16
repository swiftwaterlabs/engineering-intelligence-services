package models

import "time"

type RepositoryProcessingOptions struct {
	IncludeDetails      bool
	IncludeOwners       bool
	IncludePullRequests bool
	IncludeBranchRules  bool
	IncludeWebhooks     bool
	Organizations       []string
	Since               *time.Time
}

type TestResultProcessingOptions struct {
	Since    *time.Time
	Projects []string
}
