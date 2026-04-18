package provider

import "strings"

// Provider holds the CLI name and flag names needed to create a PR/MR.
type Provider struct {
	Name             string   // "github" | "gitlab" | "unknown"
	CLI              string   // "gh" | "glab"
	SubCmd           []string // subcommand tokens, e.g. ["pr","create"] or ["mr","create"]
	BaseBranchFlag   string   // "--base" | "--target-branch"
	BodyFlag         string   // "--body" | "--description"
	SourceBranchFlag string   // "" (implicit) | "--source-branch"
}

var github = Provider{
	Name:             "github",
	CLI:              "gh",
	SubCmd:           []string{"pr", "create"},
	BaseBranchFlag:   "--base",
	BodyFlag:         "--body",
	SourceBranchFlag: "",
}

var gitlab = Provider{
	Name:             "gitlab",
	CLI:              "glab",
	SubCmd:           []string{"mr", "create"},
	BaseBranchFlag:   "--target-branch",
	BodyFlag:         "--description",
	SourceBranchFlag: "--source-branch",
}

// Detect returns the Provider for the given remote URL.
// Falls back to GitHub when the provider cannot be determined.
func Detect(remoteURL string) Provider {
	u := strings.ToLower(remoteURL)
	switch {
	case strings.Contains(u, "gitlab"):
		return gitlab
	default:
		return github
	}
}
