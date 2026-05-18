package provider

import "strings"

// Provider holds the CLI name and commands for the detected Git host.
type Provider struct {
	Name          string   // "github" | "gitlab" | "unknown"
	CLI           string   // "gh" | "glab"
	BrowseRepoCmd []string // ["browse"] | ["repo", "view", "--web"]
	ViewPRCmd     []string // ["pr", "view", "--web"] | ["mr", "view", "--web"]
}

var github = Provider{
	Name:          "github",
	CLI:           "gh",
	BrowseRepoCmd: []string{"browse"},
	ViewPRCmd:     []string{"pr", "view", "--web"},
}

var gitlab = Provider{
	Name:          "gitlab",
	CLI:           "glab",
	BrowseRepoCmd: []string{"repo", "view", "--web"},
	ViewPRCmd:     []string{"mr", "view", "--web"},
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
