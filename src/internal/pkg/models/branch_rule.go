package models

type BranchProtectionRule struct {
	Id                           string
	Type                         string
	Repository                   *Repository
	Branch                       string
	AllowForcePush               bool
	RequirePullRequest           bool
	RequirePullRequestApprovals  bool
	RequiredPullRequestApprovers int
	IncludeAdministrators        bool
	RawData                      interface{}
}
