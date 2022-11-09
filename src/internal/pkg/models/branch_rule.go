package models

type BranchProtectionRule struct {
	Id                           string
	Repository                   *Repository
	Branch                       string
	AllowForcePush               bool
	RequirePullRequest           bool
	RequirePullRequestApprovals  bool
	RequiredPullRequestApprovers int
	IncludeAdministrators        bool
	RawData                      interface{}
}
